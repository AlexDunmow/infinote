import gql from "graphql-tag"
import { createContainer } from "unstated-next"
import * as React from "react"
import { OnboardingInput, Prospect, UserError } from "../types/types"
import { useMutation } from "@apollo/react-hooks"
import { ApolloError } from "apollo-client"

const MUTATION_ONBOARD_START = gql`
	mutation onboardStart($email: String!) {
		onboardStart(email: $email) {
			userErrors {
				message
				field
			}
			onboard {
				prospect {
					id
					email
					name
				}
			}
		}
	}
`

const MUTATION_ONBOARD_UPDATE = gql`
	mutation onboardUpdate($id: String!, $input: OnboardingInput!) {
		onboardUpdate(id: $id, input: $input) {
			userErrors {
				message
				field
			}
			onboard {
				prospect {
					id
					email
					name
				}
			}
		}
	}
`

const MUTATION_ONBOARD_FINISH = gql`
	mutation onboardFinish($id: String!) {
		onboardFinish(id: $id) {
			userErrors {
				message
				field
			}
			onboard {
				prospect {
					id
					email
					name
					onboardingComplete
				}
			}
		}
	}
`

const useOnboarding = () => {
	const [submitError, setSubmitError] = React.useState<UserError[]>([])
	const [apolloError, setApolloError] = React.useState<ApolloError | undefined>(undefined)
	const [prospect, setProspect] = React.useState<Prospect>({
		id: "",
		email: "",
		name: "",
		onboardingComplete: false
	})
	const [loading, setLoading] = React.useState<boolean>(false)
	const [current, setCurrent] = React.useState<number>(0)

	const useStartOnboarding = () => {
		const [startOnboarding, { data, loading, error }] = useMutation(MUTATION_ONBOARD_START, {
			variables: { email: prospect.email }
		})

		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)
			if (data && data.onboardStart.userErrors.length > 0) {
				setSubmitError(data.onboardStart.userErrors)
				return
			}
			if (data) {
				setProspect({ ...prospect, email: data.onboardStart.onboard.prospect.email, id: data.onboardStart.onboard.prospect.id })
				stepForward()
			}
			return
		}, [data, loading, error])

		return { startOnboarding }
	}

	const useUpdateOnboarding = () => {
		const input: OnboardingInput = {
			email: prospect.email,
			name: prospect.name
		}
		const [updateOnboarding, { data, loading, error }] = useMutation(MUTATION_ONBOARD_UPDATE, {
			variables: { id: prospect.id, input: input }
		})

		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)
			if (data && data.onboardUpdate.userErrors.length > 0) {
				setSubmitError(data.onboardUpdate.userErrors)
				return
			}
			if (data) {
				setProspect(data.onboardUpdate.onboard.prospect)
				stepForward()
			}
		}, [data, loading, error])
		return { updateOnboarding }
	}

	const useFinishOnboarding = () => {
		const [finishOnboarding, { data, loading, error }] = useMutation(MUTATION_ONBOARD_FINISH, {
			variables: { id: prospect.id }
		})

		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)
			if (data && data.onboardFinish.userErrors.length > 0) {
				setSubmitError(data.onboardFinish.userErrors)
				return
			}
			if (data) {
				stepForward()
			}
		}, [data, loading, error])

		return { finishOnboarding }
	}

	const stepForward = () => {
		setCurrent(current + 1)
	}

	const stepBack = () => {
		setCurrent(current - 1)
	}

	return {
		useStartOnboarding,
		useUpdateOnboarding,
		useFinishOnboarding,
		stepBack,
		current,
		loading,
		submitError,
		apolloError,
		prospect,
		setProspect
	}
}

export const Onboarding = createContainer(useOnboarding)
