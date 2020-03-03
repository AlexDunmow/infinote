import * as React from "react"
import { EditorContentManager } from "@convergencelabs/monaco-collab-ext"
import { Note, NoteChange, NoteEvent } from "../../types/types"

interface Props {
	sessionID: string
	contentManager: EditorContentManager
	event?: { NoteEvent: NoteEvent }
	setHistory(history: string[]): void
	history: string[]
}

const Changes = ({ sessionID, contentManager, history, event, setHistory }: Props) => {
	if (!event) {
		return <></>
	}
	const { insert, eventID, replace } = event.NoteEvent
	const arrayCheck = history.indexOf(eventID)
	if (arrayCheck === -1) {
		const newHistory = [...history, eventID]
		setHistory(newHistory)

		if (event.NoteEvent.sessionID !== sessionID) {
			if (insert) {
				contentManager.insert(insert.index, insert.text)
			}
			if (replace) {
				contentManager.replace(5, 10, "some text")
			}
		}
	}

	return <></>
}

export default Changes
