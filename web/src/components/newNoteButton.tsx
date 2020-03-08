import * as React from "react"
import { Button, SHAPE } from "baseui/button"
import { ReactNode } from "react"
import { SIZE } from "baseui/input"

interface Props {
	children: ReactNode
	isLarge?: boolean
}

const NewNoteButton = ({ children, isLarge }: Props) => {
	return (
		<Button
			shape={SHAPE.round}
			size={isLarge ? SIZE.large : SIZE.compact}
			onClick={() => alert("click")}
			overrides={{
				BaseButton: {
					style: ({ $theme }) => {
						return {
							margin: "5px",
							backgroundColor: $theme.colors.warning200
						}
					}
				}
			}}>
			{children}
		</Button>
	)
}

export default NewNoteButton
