import * as React from "react"
import { useStyletron } from "baseui"
import { Accordion, Panel } from "baseui/accordion"
import { H2 } from "baseui/typography"
import { UserContainer } from "../controllers/user"
import { ChangePassword } from "../components/changePassword"
import { ChangeSettings } from "../components/changeDetails"

interface Props {
	setShowLogin: (showLogin: boolean) => void
	showLogin: boolean
}

export const Settings = () => {
	const [css, theme] = useStyletron()
	const container: string = css({
		margin: "2rem",
	})

	const { clearErrors } = UserContainer.useContainer()

	return (
		<div className={container}>
			<H2>User Settings</H2>
			<Accordion
				onChange={({ expanded }) => {
					clearErrors()
				}}
			>
				<Panel title="Your Details">
					<ChangeSettings />
				</Panel>
				<Panel title="Panel 2">Content 2</Panel>
				<Panel title="Change Password">
					<ChangePassword />
				</Panel>
			</Accordion>
		</div>
	)
}
