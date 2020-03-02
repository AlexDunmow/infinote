import gql from "graphql-tag"
import { createContainer } from "unstated-next"
import * as React from "react"
import { UserError, User } from "../types/types"
import { useMutation, useQuery } from "@apollo/react-hooks"
import { ApolloError } from "apollo-client"

const MUTATION_CHANGE_PASSWORD = gql`
	mutation changePassword($oldPassword: String!, $password: String!) {
		changePassword(oldPassword: $oldPassword, password: $password) {
			userErrors {
				message
				field
			}
			success
		}
	}
`

const MUTATION_CHANGE_DETAILS = gql`
	mutation changeDetails($name: String!) {
		changeDetails(name: $name) {
			userErrors {
				message
				field
			}
			viewer {
				me {
					id
					name
					email
					verified
				}
			}
		}
	}
`

export const USER_QUERY = gql`
	{
		me {
			id
			name
			email
			verified
		}
	}
`

const useUser = () => {
	const [userErrors, setUserErrors] = React.useState<UserError[]>([])
	const [apolloError, setApolloError] = React.useState<ApolloError | undefined>(undefined)
	const [loading, setLoading] = React.useState<boolean>(false)
	const [user, setUser] = React.useState<User | undefined>(undefined)

	const clearErrors = () => {
		setUserErrors([])
		setApolloError(undefined)
	}

	const useChangePassword = (oldPassword: string, password: string) => {
		const [changePassword, { data, loading, error }] = useMutation(MUTATION_CHANGE_PASSWORD, {
			variables: { oldPassword, password }
		})
		const [changeSuccess, setChangeSuccess] = React.useState<boolean>(false)

		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)
			if (data && data.changePassword.userErrors.length > 0) {
				setUserErrors(data.changePassword.userErrors)
				return
			}
			if (data) {
				setChangeSuccess(data.changePassword.success)
			}
			return
		}, [data, loading, error])

		return { changePassword, changeSuccess, setChangeSuccess }
	}

	const useFetchUser = () => {
		const { data, loading, error } = useQuery(USER_QUERY)
		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)
			if (data) {
				if (data.viewer && data.viewer.me) {
					setUser(data.viewer.me)
				}
			}
			return
		}, [data, loading, error])
	}

	const useChangeDetails = (name: string) => {
		const [changeDetails, { data, loading, error }] = useMutation(MUTATION_CHANGE_DETAILS, {
			variables: { name }
		})
		const [changeSuccess, setChangeSuccess] = React.useState<boolean>(false)

		React.useEffect(() => {
			setLoading(loading)
			setApolloError(error)
			if (data && data.changeDetails.userErrors.length > 0) {
				setUserErrors(data.changeDetails.userErrors)
				return
			}
			if (data) {
				setChangeSuccess(true)
				setUser(data.changeDetails.viewer.me)
			}
			return
		}, [data, loading, error])

		return { changeDetails, changeSuccess, setChangeSuccess }
	}

	return {
		loading,
		userErrors,
		apolloError,
		useChangePassword,
		user,
		useFetchUser,
		clearErrors,
		useChangeDetails
	}
}

export const UserContainer = createContainer(useUser)
