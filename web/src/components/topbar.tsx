import * as React from "react"
import { useStyletron } from "baseui"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { StatefulTooltip } from "baseui/tooltip"
import { PLACEMENT, StatefulPopover } from "baseui/popover"
import { StatefulMenu } from "baseui/menu"
import { H6 } from "baseui/typography"
import { UserContainer } from "../controllers/user"
import { AuthContainer } from "../controllers/auth"
import { useHistory } from "react-router-dom"

export const TopBar = () => {
	const [css, theme] = useStyletron()

	// Note: get real has message value
	const hasMessage = true

	const containerStyle: string = css({
		display: "flex",
		justifyContent: "flex-end",
		alignItems: "center",
		width: "100%",
		height: "70px",
		borderBottom: `1px solid #D9D9D9`,
		fontSize: "1.5rem"
	})

	const topBarContentStyle: string = css({
		marginRight: "25px",
		display: "flex",
		alignItems: "center"
	})

	return (
		<div className={containerStyle}>
			<div className={topBarContentStyle}>
				<Messages hasMessage={hasMessage} />
				<Account />
			</div>
		</div>
	)
}

const Account = () => {
	const [css, theme] = useStyletron()
	const [isOpen, setIsOpen] = React.useState<boolean>(false)
	const { logout } = AuthContainer.useContainer()
	const { user } = UserContainer.useContainer()

	const history = useHistory()

	const userStyle: string = css({
		cursor: "pointer",
		display: "flex",
		alignItems: "center"
	})

	const userNameStyle: string = css({
		marginLeft: "0.5rem",
		marginRight: "0.5rem"
	})

	if (!user) {
		console.error("Can't load top bar. 'user' data not available.")
		return null
	}

	const menuContent = (close: () => void) => (
		<StatefulMenu
			overrides={{
				List: {
					style: ({ $theme }) => {
						return {
							outline: "none",
							minWidth: "200px"
						}
					}
				}
			}}
			items={[{ label: "Settings" }, { label: "Log Out" }]}
			onItemSelect={selection => {
				switch (selection.item.label) {
					case "Settings":
						history.push("/portal/settings")
						close()
						return
					case "Log Out":
						logout()
						return
					default:
						return
				}
			}}
		/>
	)

	return (
		<StatefulPopover
			content={({ close }) => menuContent(close)}
			onOpen={() => setIsOpen(true)}
			onClose={() => setIsOpen(false)}
			placement={PLACEMENT.bottomRight}>
			<div className={userStyle}>
				<FontAwesomeIcon icon={["fal", "user-circle"]} />
				<div className={userNameStyle}>
					<H6
						overrides={{
							Block: {
								style: {
									marginTop: "0px",
									marginBottom: "0px"
								}
							}
						}}>
						{`${user.name}`}
					</H6>
				</div>
				<FontAwesomeIcon icon={["fal", isOpen ? "chevron-up" : "chevron-down"]} size={"xs"} />
			</div>
		</StatefulPopover>
	)
}

interface IMessageProps {
	hasMessage: boolean
}

const Messages = (props: IMessageProps) => {
	const { hasMessage } = props

	const [css, theme] = useStyletron()
	const messageStyle: string = css({
		cursor: "pointer",
		marginRight: "1rem",
		display: "flex",
		alignItems: "center",
		position: "relative"
	})

	const notificationCountStyle: string = css({
		backgroundColor: theme.colors.accent,
		borderRadius: "50%",
		width: "0.7rem",
		height: "0.7rem",
		position: "absolute",
		top: "-0.1rem",
		right: "-0.2rem"
	})

	const handleClick = () => {
		console.log("not implemented")
	}

	if (hasMessage) {
		return (
			<StatefulTooltip
				showArrow
				content={() => {
					return "You have new messages"
				}}>
				<div className={messageStyle} onClick={handleClick}>
					<div className={notificationCountStyle} />
					<FontAwesomeIcon icon={["fal", "envelope"]} />
				</div>
			</StatefulTooltip>
		)
	}
	return (
		<div className={messageStyle}>
			<FontAwesomeIcon icon={["fal", "envelope"]} />
		</div>
	)
}
