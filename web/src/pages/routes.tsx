import * as React from "react"
import { Route, Switch } from "react-router-dom"
import { useStyletron } from "baseui"
import { Home } from "./home"
import LoginPage from "./loginPage"
import { AuthContainer } from "../controllers/auth"
import { Loading } from "../components/loading"
import VerifyPage from "./verify"
import EditNote from "./editNote"

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
				<Route path={"/verify/:url"} component={VerifyPage} />
				<Route path={"/login/:url"} component={LoginPage} />
				<Route path={"/verify"} component={VerifyPage} />
				<Route path={"/login"} component={LoginPage} />
				<Route path={"/note/:noteID"} component={EditNote} />
				<Route path={"/"} component={Home} />
			</Switch>
		</div>
	)
}

export default Routes
