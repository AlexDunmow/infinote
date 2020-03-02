import * as React from "react"
import { useStyletron } from "baseui"
import { Login } from "./login"

interface Props {
	setShowLogin: (showLogin: boolean) => void
	showLogin: boolean
}

export const AnimatedLogin = (props: Props) => {
	const { setShowLogin, showLogin } = props

	const [css, theme] = useStyletron()
	const containerStyle: string = css({
		padding: "2rem",
		borderRight: "none",
		borderTop: "none",
		backgroundColor: "white",
		top: "0px",
		right: "0px",
		position: "absolute",
		transition: "transform 500ms ease",
		transform: showLogin ? "translateY(0)" : "translateY(-100%)",
		boxShadow: showLogin ? "0px 0px 10px grey" : "none",
	})

	return (
		<div className={containerStyle}>
			<Login setShowLogin={setShowLogin} />
		</div>
	)
}
