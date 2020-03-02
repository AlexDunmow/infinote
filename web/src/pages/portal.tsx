import * as React from "react"
import { useStyletron } from "baseui"
import { TopBar } from "../components/topbar"
import { SideBar } from "../components/sidebar"
import { Switch, Route, Redirect } from "react-router"
import { Support } from "./support"
import { Account } from "./account"
import { UserContainer } from "../controllers/user"
import { AuthContainer } from "../controllers/auth"
import { VerificationComplete } from "./verificationComplete"
import { Settings } from "./settings"
import { NotLoggedIn } from "./notLogedIn"
import { Loading } from "../components/loading"
import { Billing } from "./billing"
import { Products } from "./products"

export const PortalInner = () => {
	const [css, theme] = useStyletron()

	const containerStyle: string = css({
		display: "flex",
		minHeight: "100vh",
		width: "100%",
	})

	const mainStyle: string = css({
		display: "flex",
		width: "100%",
		flexDirection: "column",
	})

	const contentStyle: string = css({
		width: "100%",
		height: "100%",
	})

	const { user, loading, useFetchUser } = UserContainer.useContainer()
	const { logoutRedirect, setLogoutRedirect, showVerifyComplete } = AuthContainer.useContainer()
	useFetchUser()

	// Redirect triggered by logout
	if (logoutRedirect) {
		setLogoutRedirect(false)
		return <Redirect to={"/"} />
	}

	if (!loading && !user) {
		return <NotLoggedIn />
	}

	if (!user) {
		return <Loading />
	}

	return (
		<React.Fragment>
			<div className={containerStyle}>
				<SideBar />
				<div className={mainStyle}>
					<TopBar />
					<div className={contentStyle}>
						<Switch>
							<Route path={"/portal/support"} component={Support} />
							<Route path={"/portal/account"} component={Account} />
							<Route path={"/portal/settings"} component={Settings} />
							<Route path={"/portal/billing"} component={Billing} />
							<Route path={"/portal/products"} component={Products} />
						</Switch>
					</div>
				</div>
			</div>
			{showVerifyComplete && <VerificationComplete />}
		</React.Fragment>
	)
}

export const Portal = () => {
	return (
		<UserContainer.Provider>
			<PortalInner />
		</UserContainer.Provider>
	)
}
