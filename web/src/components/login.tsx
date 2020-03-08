import * as React from "react"
import { Button, KIND } from "baseui/button"
import { useStyletron } from "baseui"
import { Spread } from "./spread"
import { ErrorBox } from "./errorBox"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { AuthContainer } from "../controllers/auth"
import { Spinner } from "baseui/spinner"
import { Spaced } from "./spaced"
import { Redirect, useHistory } from "react-router"
import { Link } from "react-router-dom"
import { ErrorMap } from "../types/types"
import { Loading } from "./loading"

interface Props {
	setShowLogin?: (showLogin: boolean) => void
	redirect?: string
}
// must supply setShowLogin if using with the animated login as it uses this prop to know when its being hidden.
// cancel button will redirect to home when this prop is not provided.

export const Login = ({ setShowLogin, redirect }: Props) => {
	const history = useHistory()

	const [css, theme] = useStyletron()
	const forgotPasswordStyle: string = css({
		marginTop: "1rem"
	})

	const { login, apolloError, userErrors, loading, clearAuthErrors, loggedIn } = AuthContainer.useContainer()
	const [inputError, setInputError] = React.useState<ErrorMap>({})

	const [email, setEmail] = React.useState<string>("")
	const [password, setPassword] = React.useState<string>("")

	const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()

		const errors: ErrorMap = {}
		let foundError = false

		if (email == "") {
			errors["email"] = "please enter your email"
			foundError = true
		}
		if (password == "") {
			errors["password"] = "please enter your password"
			foundError = true
		}

		if (foundError) {
			setInputError(errors)
			return
		}

		setInputError({})
		login(email, password).then(() => {
			if (redirect) {
				history.push(redirect)
			}
		})
	}

	const handleCancel = () => {
		clearAuthErrors()
		setInputError({})
		setEmail("")
		setPassword("")
		if (setShowLogin) {
			setShowLogin(false)
			return
		}
		return history.push("/")
	}

	if (loading) {
		return <Loading />
	}

	return (
		<form onSubmit={onSubmit}>
			<div>
				<ErrorBox apolloError={apolloError} userErrors={userErrors} />
				<FormControl label="Email" error={inputError["email"]}>
					<Input
						key={"email"}
						error={!!inputError["email"]}
						positive={false}
						value={email}
						onChange={e => {
							setEmail(e.currentTarget.value)
						}}
						placeholder={"Your email"}
					/>
				</FormControl>
				<FormControl label="Password" error={inputError["password"]}>
					<Input
						key={"password"}
						error={!!inputError["password"]}
						positive={false}
						value={password}
						type={"password"}
						onChange={e => {
							setPassword(e.currentTarget.value)
						}}
						placeholder={"Your password"}
					/>
				</FormControl>
			</div>
			<Button type="submit">Login</Button>

			<div className={forgotPasswordStyle}>
				<Link to={"/"}>Forgot Password</Link>
			</div>
		</form>
	)
}
