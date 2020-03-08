import * as React from "react"
import { Route, Switch } from "react-router-dom"
import { useStyletron } from "baseui"
import { Home } from "./home"
import LoginPage from "./loginPage"
import { AuthContainer } from "../controllers/auth"
import { useQuery } from "@apollo/react-hooks"
import { USER_QUERY } from "../controllers/user"
import { Loading } from "../components/loading"

const Routes = () => {
	const [css, theme] = useStyletron()
	const routeStyle: string = css({
		width: "100%",
		minHeight: "100vh"
	})

	const auth = AuthContainer.useContainer()

	if (!auth.check.checked || auth.check.checking) {
		console.log("auth not checked", auth)
		return <Loading />
	}

	return (
		<div className={routeStyle}>
			<Switch>
				<Route path={"/login/:url"} component={LoginPage} />
				<Route path={"/note/:noteID"} component={Home} />
				<Route path={"/"} component={Home} />
			</Switch>
		</div>
	)
}

export default Routes
