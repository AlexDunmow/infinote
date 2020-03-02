import * as React from "react"
import { KIND as ButtonKIND } from "baseui/button"
import { Button } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { StyleObject } from "styletron-react"
import { Spaced } from "./spaced"
import { ErrorBox } from "./errorBox"
import { useStyletron } from "baseui"
import { Spread } from "./spread"
import { Onboarding } from "../controllers/onboarding"

export const OnboardSummary = () => {
	const [css, theme] = useStyletron()
	const { useFinishOnboarding, stepBack, apolloError, submitError, prospect, loading } = Onboarding.useContainer()
	const { finishOnboarding } = useFinishOnboarding()

	const handleClick = () => {
		finishOnboarding()
	}

	return (
		<div>
			<ErrorBox apolloError={apolloError} userErrors={submitError} />
			<p>Please review the information provided</p>
			<div className={css(detailsStyle)}>
				<span>
					<strong>Email:</strong> {prospect.email}
				</span>
				<span>
					<strong>First Name:</strong> {prospect.name}
				</span>
			</div>
			<Spread>
				<Button kind={ButtonKIND.secondary} onClick={stepBack}>
					Previous
				</Button>
				<Spaced>
					{loading && <Spinner />}
					<Button onClick={handleClick}>Submit</Button>
				</Spaced>
			</Spread>
		</div>
	)
}

const detailsStyle: StyleObject = {
	padding: "1em",
	display: "flex",
	flexDirection: "column"
}
