import * as React from "react"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Button, KIND } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { Spaced } from "./spaced"
import { ErrorBox } from "./errorBox"
import { Spread } from "./spread"
import { Onboarding } from "../controllers/onboarding"
import { ErrorMap, UserError } from "../types/types"

export const OnboardingContact = () => {
	const { useUpdateOnboarding, stepBack, prospect, setProspect, apolloError, submitError, loading } = Onboarding.useContainer()
	const { updateOnboarding } = useUpdateOnboarding()

	const [inputError, setInputError] = React.useState<ErrorMap>({})

	const handleClick = async () => {
		const inputErrors: ErrorMap = {}
		let foundError: boolean = false
		if (prospect.name == "") {
			inputErrors["name"] = "You must enter your first name"
			foundError = true
		}
		if (foundError) {
			setInputError(inputErrors)
			return
		}
		updateOnboarding()
	}

	return (
		<div>
			<ErrorBox apolloError={apolloError} userErrors={submitError} />
			<p>We need some of your basic details.</p>
			<FormControl label="First Name" error={inputError["name"]}>
				<Input
					key={"name"}
					error={!!inputError["name"]}
					positive={false}
					value={prospect.name}
					onChange={e => {
						setProspect({ ...prospect, name: e.currentTarget.value })
					}}
					placeholder={"Your first name"}
				/>
			</FormControl>
			<Spread>
				<Button kind={KIND.secondary} onClick={stepBack}>
					Back
				</Button>
				<Spaced>
					{loading && <Spinner />}
					<Button onClick={handleClick}>Next</Button>
				</Spaced>
			</Spread>
		</div>
	)
}
