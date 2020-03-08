import * as React from "react"
import { gql } from "apollo-boost"
import { Redirect, useParams } from "react-router-dom"
import { useQuery } from "@apollo/react-hooks"
import { INote } from "../types/types"
import { Loading } from "../components/loading"
import { AuthContainer } from "../controllers/auth"
import NoteEditor from "../components/Editor/editor"

const GETNOTE = gql`
	query getNote($noteID: ID!) {
		noteByID(noteID: $noteID) {
			id
			name
			body
			done
			owner {
				id
			}
		}
	}
`

interface Props {
	noteID: string
}

const EditNote = () => {
	const { noteID } = useParams<{ noteID: string }>()

	const auth = AuthContainer.useContainer()

	const { loading, error, data } = useQuery<{ noteByID: INote }>(GETNOTE, {
		variables: { noteID: noteID || "" }
	})

	if (loading || !data) {
		console.log(
			"auth.check.checked:",
			auth.check.checked,
			"auth.check.checking",
			auth.check.checking,
			"loading:",
			loading,
			"auth.loading:",
			auth.loading,
			"data:",
			data,
			"error:",
			error,
			"logged in:",
			auth.loggedIn
		)

		// auth.check.checked: true auth.check.checking false loading: true auth.loading: false data: {} error: undefined logged in: false
		// home.tsx?d54c:64 auth.check.checked: true auth.check.checking false loading: true auth.loading: false data: {} error: undefined logged in: false

		if (auth.check.checked && !auth.loggedIn && !auth.loading) {
			console.log("redirecting?")
			const url = btoa(window.location.pathname)

			return <Redirect to={`/login/${url}`} />
		} else if (auth.check.checked && auth.user && !auth.user.verified) {
			console.log("redirecting to verification?")
			const url = btoa(window.location.pathname)

			return <Redirect to={`/verify/${url}`} />
		}

		return (
			<div>
				Loading Note...
				<Loading />
			</div>
		)
	}

	return <NoteEditor note={data.noteByID} />
}

export default EditNote
