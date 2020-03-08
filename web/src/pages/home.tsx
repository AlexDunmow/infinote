import * as React from "react"
import { RouteComponentProps, Redirect } from "react-router-dom"
import { StyleObject } from "styletron-react"
import { Button, SHAPE } from "baseui/button"
import { Spaced } from "../components/spaced"
import { AnimatedLogin } from "../components/animatedLogin"
import { useParams, useLocation, useHistory, useRouteMatch } from "react-router-dom"
import { useStyletron } from "baseui"
import { AuthContainer } from "../controllers/auth"
import { Loading } from "../components/loading"
import Notes from "../components/notes"
import NoteEditor from "../components/Editor/editor"
import { useQuery } from "@apollo/react-hooks"
import { gql } from "apollo-boost"
import { INote } from "../types/types"
import { ButtonGroup } from "baseui/button-group"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { SIZE } from "baseui/input"
import { faPlus } from "@fortawesome/pro-regular-svg-icons"
import NewNoteButton from "../components/newNoteButton"
interface IProps extends RouteComponentProps {}

const Logo = require("../assets/images/Ninja-Software-Hero.png")

export const Home = (props: IProps) => {
	const [css, theme] = useStyletron()
	const auth = AuthContainer.useContainer()
	const background = css({
		minHeight: "100vh",
		maxHeight: "100vh",
		width: "100%",
		backgroundImage: `url(${Logo})`,
		backgroundRepeat: "repeat",
		backgroundSize: "20%",
		display: "flex",
		justifyContent: "center",
		alignItems: "center"
	})

	const controlCls = css({
		height: "100px",
		alignItems: "center",
		justifyContent: "center",
		display: "flex",
		width: "100%"
	})

	const notesCls = css({
		flex: 1,
		display: "flex"
	})

	const container = css({
		background: "white",
		maxWidth: "1200px",
		maxHeight: "100vh",
		width: "100%",
		margin: "0 auto",
		height: "100vh",
		display: "flex",
		flexDirection: "column"
	})

	if (!auth.loggedIn) {
		console.log("not logged in", auth.check)
		return <Redirect to={"/login"} />
	}

	return (
		<div className={background}>
			<div className={container}>
				<Notes />
				<div className={controlCls}>
					<NewNoteButton>
						<FontAwesomeIcon icon={"user-plus"} />
					</NewNoteButton>
					<NewNoteButton isLarge={true}>
						<FontAwesomeIcon icon={"plus"} />
					</NewNoteButton>
					<NewNoteButton>
						<FontAwesomeIcon icon={"user-plus"} />
					</NewNoteButton>
				</div>
			</div>
		</div>
	)
}
