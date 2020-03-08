import * as React from "react"
import { ApolloLink, Operation } from "apollo-link"
import { WebSocketLink } from "apollo-link-ws"
import { SubscriptionClient } from "subscriptions-transport-ws"
import { createUploadLink } from "apollo-upload-client"
import { onError } from "apollo-link-error"
import { DefinitionNode } from "graphql"

/**
 * Additional HTTP Error Handling outside of GraphQL Response
 */
const errorHttpLink = onError(({ networkError }) => {
	// Fix up JSON parse error...
	if (networkError && "statusCode" in networkError && networkError.statusCode == 413) {
		networkError.message = "413 Request Entity Too Large"
	}
})

/**
 * Apollo Link Splitter context for accessing the auth features.
 */
export const ApolloLinkSplitterContext = React.createContext<ApolloLinkSplitter | null>(null)

/**
 * Apollo Link Splitter to handle auth and file upload.
 */
export class ApolloLinkSplitter {
	authToken: string
	link: ApolloLink
	wsClient: SubscriptionClient

	/**
	 * Setup Apollo Link Splitter.
	 *
	 * @param endPoint
	 * @param isSecure
	 * @param authToken
	 */
	constructor(endPoint: string, isSecure: boolean, authToken?: string) {
		let httpEndpoint: string
		let wsEndpoint: string

		if (isSecure) {
			httpEndpoint = "https://" + endPoint
			wsEndpoint = "wss://" + endPoint
		} else {
			httpEndpoint = "http://" + endPoint
			wsEndpoint = "ws://" + endPoint
		}

		this.authToken = authToken ? authToken : ""

		// Setup Apollo HTTP Link with Uploading Support
		const uploadLink = ApolloLink.from([errorHttpLink, createUploadLink({ uri: httpEndpoint })])

		// Setup Apollo Split Link
		this.wsClient = new SubscriptionClient(wsEndpoint, {
			reconnect: true,
			connectionParams: () => {
				return {
					authorization: this.authToken
				}
			}
		})

		this.link = ApolloLink.split(
			(op: Operation) => {
				// Define list of operations involving file upload here
				return op.query.definitions.some(
					(definition: DefinitionNode) =>
						definition.kind === "OperationDefinition" &&
						definition.operation === "mutation" &&
						definition.name &&
						definition.name.value === "onboardingFileUpload"
				)
			},

			uploadLink,
			new WebSocketLink(this.wsClient)
		)
	}

	/**
	 * Update or Unset the Authorization Token and restart the Websocket
	 *
	 * @param token
	 *     The token to set; empty string to unset it.
	 */
	updateAuthToken = (token: string) => {
		this.authToken = token

		// Close connection and let Websocket Reconnect
		this.wsClient.close(true, true)
	}

	/**
	 * Reset Websocket Link.
	 */
	resetLink = () => this.wsClient.close(true, true)

	/**
	 * Return a Apollo Link with Splitting features.
	 *
	 * @returns
	 *     An ApolloLink object for ApolloClient
	 */
	getLink = () => this.link
}
