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
}

export const EmailVerify = (props: IProps) => {
	const { code } = props

	const [css, theme] = useStyletron()
	const cardContainer: string = css({
		height: "100vh",
		display: "flex",
		alignItems: "center",
		justifyContent: "center",
	})

	const { userErrors, loading, loggedIn, verify } = AuthContainer.useContainer()

	const [verifyCode, setVerifyCode] = React.useState<string>(code ? code : "")
	const [inputError, setInputError] = React.useState<string>("")

	const handleClick = () => {
		// Note: validation?
		verify(verifyCode)
	}

	if (loggedIn) {
		return <Redirect to={"/portal"} />
	}

	return (
		<div className={cardContainer}>
			<Card overrides={{ Root: { style: { flexGrow: 0.3 } } }}>
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
		</div>
	)
}
