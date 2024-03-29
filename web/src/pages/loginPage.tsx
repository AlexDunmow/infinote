import * as React from "react"
import { styled, useStyletron } from "baseui"
import { Login } from "../components/login"
import { useQuery } from "@apollo/react-hooks"
import { USER_QUERY } from "../controllers/user"
import { Loading } from "../components/loading"
import { AuthContainer } from "../controllers/auth"
import { Redirect } from "react-router"
import { useParams } from "react-router-dom"

const LogoPath = require("../assets/images/logo.svg")

const LoginPage = () => {
	const { loading, check, loggedIn } = AuthContainer.useContainer()
	if (loggedIn) {
		return <Redirect to={"/"} />
	}

	const [css, theme] = useStyletron()
	const Container = styled("div", {
		width: "100vw",
		height: "100vh",
		background: `url(https://source.unsplash.com/random/${window.innerWidth}x${window.innerHeight})`
	})

	const FormContainer = styled("div", {
		maxWidth: "700px",
		maxHeight: "700px",
		padding: "50px",
		position: "fixed",
		left: "50px",
		bottom: "50px",
		background: theme.colors.background
	})

	const Logo = styled("div", {
		textAlign: "center"
	})

	const logoCss = css({
		maxWidth: "90px"
	})

	const { url } = useParams<{ url: string }>()
	const redirect = atob(url)

	return (
		<Container>
			<FormContainer>
				{!check.checked && <Loading />}
				{check.checked && (
					<Logo>
						<img className={logoCss} src={LogoPath} />
						<h1>Infinote</h1>
					</Logo>
				)}
				<Login redirect={redirect} />
			</FormContainer>
		</Container>
	)
}

export default LoginPage
