import * as React from "react"
import { useStyletron } from "baseui"
import { StyleObject } from "styletron-standard"

interface IOverrides {
	container?: StyleObject
}

interface IProps {
	children: (JSX.Element | boolean)[]
	overrides?: IOverrides
}

// Used to spread 2 or more components by using flex and space-between. eg. Use to place 2 buttons at opposing ends of a container
export const Spread = (props: IProps) => {
	const [useCss, theme] = useStyletron()
	const containerOverrides: StyleObject | undefined = props.overrides ? props.overrides.container : undefined
	const container = useCss({
		display: "flex",
		justifyContent: "space-between",
		alignItems: "center",
		...containerOverrides,
	})
	return <div className={container}>{props.children}</div>
}
