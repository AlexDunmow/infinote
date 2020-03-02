import * as React from "react"
import { RouteComponentProps, Redirect } from "react-router-dom"
import { StyleObject } from "styletron-react"
import { Button } from "baseui/button"
import { Spaced } from "../components/spaced"
import { AnimatedLogin } from "../components/animatedLogin"
import { useParams, useLocation, useHistory, useRouteMatch } from "react-router-dom"
import { useStyletron } from "baseui"
import { AuthContainer } from "../controllers/auth"
import { Loading } from "../components/loading"
import Notes from "../components/notes"
import NoteEditor from "../components/editor"
import { useQuery } from "@apollo/react-hooks"
import { gql } from "apollo-boost"
import { Note } from "../types/types"

interface IProps extends RouteComponentProps {}

const Logo = require("../assets/images/Ninja-Software-Hero.png")

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

export const Home = (props: IProps) => {
	const { noteID } = useParams<{ noteID: string }>()

	const { loading, error, data } = useQuery<{ noteByID: Note }>(GETNOTE, {
		variables: { noteID: noteID || "" }
	})

	const [css, theme] = useStyletron()
	const { loggedIn, checked } = AuthContainer.useContainer()
	const background = css({
		minHeight: "100vh",
		width: "100%",
		backgroundImage: `url(${Logo})`,
		backgroundRepeat: "repeat",
		backgroundSize: "20%",
		display: "flex",
		justifyContent: "center",
		alignItems: "center"
	})

	const container = css({
		background: "white",
		maxWidth: "1200px",
		width: "100%",
		margin: "0 auto",
		height: "100vh"
	})

	if (!checked || loading || !data) {
		return <Loading />
	}
	if (!loggedIn) {
		return <Redirect to={"/login"} />
	}

	console.log("note: ", data)

	return (
		<div className={background}>
			<div className={container}>
				<Notes />
				<NoteEditor note={data.noteByID} />
			</div>
		</div>
	)
}
