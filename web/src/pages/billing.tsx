import * as React from "react"
import { useStyletron } from "baseui"
import { StyleObject } from "styletron-react"

export const Billing = () => {
	const [css, theme] = useStyletron()

	const containerStyle: StyleObject = {
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		height: "100%",
	}

	return <div className={css(containerStyle)}>Billing stuff goes here</div>
}
