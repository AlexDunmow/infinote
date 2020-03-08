import * as React from "react"
import { gql } from "apollo-boost"
import { useQuery } from "@apollo/react-hooks"
import { INote } from "../types/types"
import { Loading } from "./loading"
import * as ReactMarkdown from "react-markdown"
import { Card, StyledAction, StyledBody } from "baseui/card"
import { styled, useStyletron } from "baseui"
import { ButtonGroup } from "baseui/button-group"
import { Button, SIZE, SHAPE } from "baseui/button"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { Display4 } from "baseui/typography"

interface Props {
	note: INote
}

const NoteCard = styled("div", ({ $theme }) => {
	return {
		backgroundColor: $theme.colors.backgroundSecondary,
		maxWidth: "600px",
		borderRadius: "5px",
		margin: "5px"
	}
})

const NoteContents = styled("div", ({ $theme }) => {
	return {
		padding: "20px"
	}
})

const Note = ({ note }: Props) => {
	return (
		<NoteCard>
			<NoteContents>
				<Display4 marginBottom="scale500">{note.name}</Display4>
				<ReactMarkdown source={note.body} />
			</NoteContents>
			<div>
				<Button size={SIZE.mini} shape={SHAPE.pill}>
					<FontAwesomeIcon icon={"edit"} />
				</Button>
				<Button size={SIZE.mini} shape={SHAPE.pill}>
					<FontAwesomeIcon icon={"trash"} />
				</Button>
			</div>
		</NoteCard>
	)
}

export default Note
