import * as React from "react"
import { useStyletron } from "baseui"
import { StyleObject } from "styletron-react"

interface IOverrides {
	container?: StyleObject
	space?: StyleObject
}

interface IProps {
	children: (JSX.Element | boolean)[]
	overrides?: IOverrides
}

// Applies a right margin to child components. eg. Use to place a space between multiple button components.
export const Spaced = (props: IProps) => {
	const [useCss, theme] = useStyletron()
	const containerOverrides: StyleObject | undefined = props.overrides ? props.overrides.container : undefined
	const container = useCss({
		display: "flex",
		alignItems: "center",
		...containerOverrides,
	})
	const spaceOverrides: StyleObject | undefined = props.overrides ? props.overrides.space : undefined
	const space = useCss({
		marginRight: "0.5rem",
		display: "flex",
		alignItems: "center",
		...spaceOverrides,
	})
	return (
		<div className={container}>
			{props.children.map((element, index) => {
				if (element) {
					return (
						<div key={"spaced-" + index} className={space}>
							{element}
						</div>
					)
				}
				return null
			})}
		</div>
	)
}
