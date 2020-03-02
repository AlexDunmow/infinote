import * as React from "react"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Button, KIND } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { Spaced } from "./spaced"
import { Spread } from "./spread"
import { ErrorBox } from "./errorBox"
import { Onboarding } from "../controllers/onboarding"

interface IProps {
	goHome: () => void
}

export const OnboardingStart = (props: IProps) => {
	const { goHome } = props
	const { useStartOnboarding, useUpdateOnboarding, loading, submitError, apolloError, prospect, setProspect } = Onboarding.useContainer()
	const { startOnboarding } = useStartOnboarding()
	const { updateOnboarding } = useUpdateOnboarding()

	const [inputError, setInputError] = React.useState<string>("")

	const onSubmit = async () => {
		setInputError("")
		const email = prospect.email
		if (!email || email == "" || email.indexOf("@") < 0) {
			setInputError("Please enter a valid email")
			return
		}
		if (prospect.id) {
			updateOnboarding()
			return
		}
		startOnboarding()
	}

	return (
		<div>
			<ErrorBox apolloError={apolloError} userErrors={submitError} />
			<p>Lets get started by entering your email address.</p>
			<FormControl label="Email" error={inputError}>
				<Input
					error={!!inputError}
					positive={false}
					value={prospect.email}
					onChange={e => {
						setProspect({ ...prospect, email: e.currentTarget.value })
					}}
					placeholder={"Your email address"}
				/>
			</FormControl>
			<Spread>
				<Button kind={KIND.secondary} onClick={goHome}>
					Return to homepage
				</Button>
				<Spaced>
					{loading && <Spinner />}
					<Button onClick={onSubmit}>Next</Button>
				</Spaced>
			</Spread>
		</div>
	)
}
