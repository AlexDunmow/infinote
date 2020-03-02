import * as React from "react"
import { OnboardingStart } from "../components/onboardingStart"
import { OnboardingContact } from "../components/onboardingContact"
import { OnboardSummary } from "../components/onboardSummary"
import { HeadingLevel, Heading } from "baseui/heading"
import { Card } from "baseui/card"
import { Button, KIND } from "baseui/button"
import { RouteComponentProps } from "react-router"
import { useStyletron } from "baseui"
import { StyleObject } from "styletron-react"
import { Onboarding as OnboardingContainer } from "../controllers/onboarding"
import { ProgressSteps } from "../components/progressSteps"

interface Props extends RouteComponentProps {}

const OnboardingInner = (props: Props) => {
	const [css, theme] = useStyletron()
	const { current, prospect } = OnboardingContainer.useContainer()

	const goHome = () => {
		props.history.push("/")
	}

	return (
		<div>
			<div className={css(cardContainer)}>
				<Card
					overrides={{
						Root: {
							style: {
								flexGrow: 1,
								minWidth: "350px",
								maxWidth: "40%",
								backgroundColor: "rgba(255,255,255,0.8)",
								borderTopLeftRadius: "5px",
								borderTopRightRadius: "5px",
								borderBottomLeftRadius: "5px",
								borderBottomRightRadius: "5px",
								backdropFilter: "blur(10px)",
							},
						},
					}}
				>
					<HeadingLevel>
						<Heading>Create Account</Heading>
						<div>
							<ProgressSteps current={current}>
								<div>
									<OnboardingStart goHome={goHome} />
								</div>
								<div>
									<OnboardingContact />
								</div>
								<div>
									<OnboardSummary />
								</div>
								<div>
									<p>Nearly There! We've sent a confirmation email to {prospect.email}</p>
									<p>Please follow the instructions to complete your account setup.</p>
									<p>Stay Wholesome!</p>
									<Button kind={KIND.secondary} onClick={goHome}>
										Home
									</Button>
								</div>
							</ProgressSteps>
						</div>
					</HeadingLevel>
				</Card>
			</div>
		</div>
	)
}

export const Onboarding = (props: Props) => {
	return (
		<OnboardingContainer.Provider>
			<OnboardingInner {...props} />
		</OnboardingContainer.Provider>
	)
}

const cardContainer: StyleObject = {
	height: "100vh",
	display: "flex",
	alignItems: "center",
	justifyContent: "center",
}
