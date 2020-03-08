import * as React from "react"
import { Card } from "baseui/card"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Spaced } from "../components/spaced"
import { Button } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { ErrorBox } from "../components/errorBox"
import { useStyletron } from "baseui"
import { AuthContainer } from "../controllers/auth"
import { Redirect } from "react-router"

interface IProps {
	code?: string
	redirect?: string
}

export const EmailVerify = (props: IProps) => {
	console.log("EMAIL VERIFY")

	const { code } = props

	const { userErrors, loading, user, verify } = AuthContainer.useContainer()

	const [verifyCode, setVerifyCode] = React.useState<string>(code ? code : "")
	const [inputError, setInputError] = React.useState<string>("")

	if (user && user.verified) {
		return <Redirect to={"/"} />
	}

	const handleClick = () => {
		// Note: validation?
		if (!user) {
			return
		}
		verify(verifyCode, user.email)
	}

	return (
		<Card>
			<div>
				<ErrorBox userErrors={userErrors} />
				<FormControl label="Verification Code" error={inputError}>
					<Input
						error={!!inputError}
						positive={false}
						value={verifyCode}
						onChange={e => setVerifyCode(e.currentTarget.value)}
						placeholder={"Enter your verification code"}
					/>
				</FormControl>
				<Spaced>
					<Button onClick={handleClick}>Verify</Button>
					{loading && <Spinner />}
				</Spaced>
			</div>
		</Card>
	)
}
