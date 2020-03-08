import { createContainer } from "unstated-next"
import * as React from "react"
import { User, UserError } from "../types/types"
import { ApolloError } from "apollo-client"
import { resetClient } from "../app"
import { useQuery } from "@apollo/react-hooks"
import { USER_QUERY } from "./user"

interface IAuthError {
	error: string
	message: string
}

const getUserError = async (response: Response): Promise<string> => {
	const defaultMessage: string = "Something went wrong. Please try again."
	try {
		const jsonResp: IAuthError = await response.json()
		if (jsonResp.message) {
			return jsonResp.message
		}
		return defaultMessage
	} catch {
		return defaultMessage
	}
}

const useAuth = () => {
	const [userErrors, setUserErrors] = React.useState<UserError[]>([])
	const [apolloError, setApolloError] = React.useState<ApolloError | undefined>(undefined)
	const [loggedIn, setLoggedIn] = React.useState<boolean>(false)
	const [logoutRedirect, setLogoutRedirect] = React.useState<boolean>(false)
	const [isLoading, setLoading] = React.useState<boolean>(true)
	const [isVerified, setVerified] = React.useState<boolean>(false)
	const [showVerifyComplete, setShowVerifyComplete] = React.useState<boolean>(false)
	const [user, setUser] = React.useState<User>()
	const [check, setChecked] = React.useState<{ checking: boolean; checked: boolean }>({ checking: true, checked: false })

	const { data, loading, error } = useQuery(USER_QUERY)

	if (!check.checked && !loading) {
		setChecked({ checked: true, checking: false })
		setLoading(false)

		if (data && data.me) {
			setLoggedIn(true)
			setUser(data.me)
			setVerified(data.me.verified)
		} else {
			setLoggedIn(false)
			setUser(undefined)
		}
	}

	const login = async (email: string, password: string) => {
		await setLoading(true)
		const response = await fetch("/api/auth/login", {
			method: "POST",
			body: JSON.stringify({ email, password })
		})
		if (response.status === 200) {
			setLoading(false)
			resetClient()
			clearAuthErrors()
			setLoggedIn(true)
			return
		}
		const userError: string = await getUserError(response)
		setUserErrors([{ message: userError, field: [] }])
		setLoading(false)
		return
	}

	const logout = async () => {
		await setLoading(true)
		const response = await fetch("/api/auth/logout", {
			method: "POST"
		})
		if (response.status === 200) {
			setLoading(false)
			resetClient()
			setLoggedIn(false)
			clearAuthErrors()
			setLogoutRedirect(true)
			return
		}
		setUserErrors([{ message: "There was a problem logging you out", field: [] }])
		setLoading(false)
		return
	}

	const verify = async (token: string, email: string) => {
		await setLoading(true)
		const response = await fetch("/api/auth/verify_account", {
			method: "POST",
			body: JSON.stringify({ token, email })
		})
		if (response.status === 200) {
			setLoading(false)
			setUser({ ...data.me, verified: true })
			setLoggedIn(true)
			setShowVerifyComplete(true)
			return
		}
		const userError: string = await getUserError(response)
		setUserErrors([{ message: userError, field: [] }])
		setLoading(false)
		return
	}

	const clearAuthErrors = () => {
		setUserErrors([])
		setApolloError(undefined)
	}

	return {
		login,
		logout,
		verify,
		loading: isLoading,
		userErrors,
		apolloError,
		clearAuthErrors,
		loggedIn,
		logoutRedirect,
		check,
		setLogoutRedirect,
		showVerifyComplete,
		user,
		setShowVerifyComplete
	}
}

export const AuthContainer = createContainer(useAuth)
