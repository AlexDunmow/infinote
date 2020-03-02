import * as React from "react"
import { Login } from "../components/login"
import { useStyletron } from "baseui"
import { H6 } from "baseui/typography"

export const NotLoggedIn = () => {
	const [css, theme] = useStyletron()
	const containerStyle: string = css({
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		height: "100vh",
	})
	const contentStyle: string = css({
		padding: "2rem",
		boxShadow: "0px 0px 10px grey",
	})
	return (
		<div className={containerStyle}>
			<div className={contentStyle}>
				<H6>You need to login to view this page.</H6>
				<Login />
			</div>
		</div>
	)
}
