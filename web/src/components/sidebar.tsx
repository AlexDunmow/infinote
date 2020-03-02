import * as React from "react"
import { useStyletron } from "baseui"
import { Link } from "react-router-dom"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { IconName } from "@fortawesome/fontawesome-svg-core"
import { H6 } from "baseui/typography"

export const SideBar = () => {
	const [css, theme] = useStyletron()

	const Logo = require("../assets/images/NinjaSoftwareLogo.svg")

	const containerStyle: string = css({
		display: "flex",
		flexDirection: "column",
		minHeight: "100vh",
		width: "180px",
		flexShrink: 0,
		backgroundImage: "linear-gradient(#FFFFFF, #000000)",
		boxShadow: "0px 3px 6px #00000029",
	})

	const logoContainerStyle: string = css({
		width: "100%",
		height: "200px",
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
	})

	const logoStyle: string = css({
		width: "80%",
	})

	return (
		<div className={containerStyle}>
			<div className={logoContainerStyle}>
				<img className={logoStyle} src={Logo} alt={"balloon"} />
			</div>
			<div>
				<SideMenuButton selected={true} icon={"hand-holding-box"} label={"Products"} url={"/portal/products"} />
				<SideMenuButton selected={false} icon={"headset"} label={"Support"} url={"/portal/support"} />
				<SideMenuButton selected={false} icon={"user"} label={"Account"} url={"/portal/account"} />
				<SideMenuButton selected={false} icon={"money-check-alt"} label={"Billing"} url={"/portal/billing"} />
			</div>
		</div>
	)
}

interface ButtonProps {
	selected: boolean
	label: string
	icon: IconName
	url: string
}

const SideMenuButton = (props: ButtonProps) => {
	const [css, theme] = useStyletron()
	const { label, icon, url } = props

	const selected = window.location.pathname.startsWith(url)

	const buttonStyle: string = css({
		height: "150px",
		backgroundColor: selected ? "rgba(0, 0, 0, 0.5)" : "transparent",
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		color: "white",
		textAlign: "center",
		":hover": {
			backgroundColor: "rgba(0, 0, 0, 0.5)",
		},
	})

	const linkStyle: string = css({
		textDecoration: "none",
	})

	return (
		<Link to={url} className={linkStyle}>
			<div className={buttonStyle}>
				<div>
					<FontAwesomeIcon icon={["fal", icon]} size={"4x"} />
					<div>
						<H6
							color={"white"}
							overrides={{
								Block: {
									style: {
										marginTop: "0.5rem",
										marginBottom: "0px",
									},
								},
							}}
						>
							{label}
						</H6>
					</div>
				</div>
			</div>
		</Link>
	)
}
