import * as React from "react"
import { gql } from "apollo-boost"
import { useQuery } from "@apollo/react-hooks"
import { INote } from "../types/types"
import { Loading } from "./loading"
import Note from "./note"
import { styled } from "baseui"

const GETNOTES = gql`
	query notes {
		notes {
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

const Container = styled("div", {
	width: "100%",
	flex: 1,
	overflowY: "auto",
	alignItems: "center",
	display: "flex",
	flexDirection: "column",
	justifyContent: "center"
})

const Notes = () => {
	const { loading, error, data } = useQuery<{ notes: INote[] }>(GETNOTES)

	if (loading || !data || !data.notes) {
		return <Loading />
	}

	console.log(data, "data")

	return (
		<Container>
			{data.notes.map(function(note, index) {
				return <Note note={note} key={note.id} />
			})}
		</Container>
	)
}

export default Notes
