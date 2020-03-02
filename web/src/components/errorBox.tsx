import * as React from "react"
import { KIND, Notification } from "baseui/notification"
import { UserError } from "../types/types"
import { ApolloError } from "apollo-client"

interface IProps {
	apolloError?: ApolloError
	userErrors: UserError[]
}

// Accepts a userError and optionally an apolloError. userError is given preference in display as it should be a friendly message
export const ErrorBox = (props: IProps) => {
	const { userErrors, apolloError } = props
	// return nothing when no error
	if (!apolloError && !userErrors) {
		return null
	}
	// if the error to be displayed is a userError, check the length is > 0 so that it doesnt display empty box
	if (!apolloError && userErrors && userErrors.length < 1) {
		return null
	}
	// Prefer to send userError if it exists as this should contain the friendly message, so attempt to return this
	// first if it is available
	if (userErrors && userErrors.length > 0) {
		return (
			<div>
				<Notification kind={KIND.warning}>
					{userErrors.map((error: UserError, index) => {
						return <div key={"error-row" + index}>{String(error.message)}</div>
					})}
				</Notification>
			</div>
		)
	}
	// Lastly attempt to return the apollo error if it is available
	if (apolloError) {
		return (
			<div>
				<Notification kind={KIND.warning}>{apolloError}</Notification>
			</div>
		)
	}
	return null
}
