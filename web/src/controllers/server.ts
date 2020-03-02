import gql from "graphql-tag"
import { createContainer } from "unstated-next"
import * as React from "react"
import { UserError } from "../types/types"
import { ServerPlan, ServerInstance } from "../types/server"
import { useMutation } from "@apollo/react-hooks"
import { ApolloError } from "apollo-client"

const QUERY_VPS_PLAN_LIST = gql`
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

const QUERY_VPS_LIST = gql`
	{
		vultrServerList {
			InstanceID
			servers {
				AllowedBandwidth
				AppID
			}
		}
	}
`

const MUTATION_VPS_CREATE = gql`
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
const MUTATION_VPS_DELETE = gql``
const MUTATION_VPS_START = gql``
const MUTATION_VPS_STOP = gql``
const MUTATION_VPS_RESTART = gql``

const useOnboarding = () => {}

export const Onboarding = createContainer(useOnboarding)
