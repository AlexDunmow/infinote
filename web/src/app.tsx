import * as React from "react"
import MetaTags from "react-meta-tags"
import { ApolloProvider, useQuery } from "@apollo/react-hooks"
import { ApolloClient, InMemoryCache } from "apollo-boost"
import { Client as Styletron } from "styletron-engine-atomic"
import { Provider as StyletronProvider } from "styletron-react"
import { BaseProvider } from "baseui"
import { LightTheme, DarkTheme } from "./themeOverrides"
import { loadIcons } from "./helpers/loadicons"
import { ApolloLinkSplitter, ApolloLinkSplitterContext } from "./apollo"
import Routes from "./pages/routes"
import { BrowserRouter as Router, Route } from "react-router-dom"
import { AuthContainer } from "./controllers/auth"

import "./app.css"
import { USER_QUERY } from "./controllers/user"

const engine = new Styletron()

loadIcons()

var loc = window.location
const GRAPHQL_ENDPOINT = "//" + loc.host + "/api/gql/query"
const apolloLinkSplitter = new ApolloLinkSplitter(GRAPHQL_ENDPOINT, loc.protocol === "https:")
const apolloClient = new ApolloClient({
	cache: new InMemoryCache(),
	link: apolloLinkSplitter.getLink()
})

export const resetClient = () => {
	apolloClient.resetStore()
	apolloLinkSplitter.resetLink()
}

const App = () => {
	const [darkTheme, setDarkTheme] = React.useState<boolean>(false)

	return (
		<StyletronProvider value={engine}>
			<BaseProvider theme={darkTheme ? DarkTheme : LightTheme}>
				<ApolloProvider client={apolloClient}>
					<ApolloLinkSplitterContext.Provider value={apolloLinkSplitter}>
						<MetaTags>
							<title>Infinote</title>
							<meta name="viewport" content="width=device-width, initial-scale=1.0" />
							<meta id="meta-description" name="description" content="Some description." />
							<meta id="og-title" property="og:title" content="MyApp" />
							<meta id="og-image" property="og:image" content="path/to/image.jpg" />
						</MetaTags>
						<AuthContainer.Provider>
							<Router>
								<Routes />
							</Router>
						</AuthContainer.Provider>
					</ApolloLinkSplitterContext.Provider>
				</ApolloProvider>
			</BaseProvider>
		</StyletronProvider>
	)
}

export { App }
