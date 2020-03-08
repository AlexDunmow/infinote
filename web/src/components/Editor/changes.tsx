import * as React from "react"
import { EditorContentManager, RemoteCursorManager } from "@convergencelabs/monaco-collab-ext"
import { CursorPlacement, INote, NoteChange, NoteEvent } from "../../types/types"
import { useState } from "react"
import { objectSize } from "../../helpers"
import { RemoteCursor } from "@convergencelabs/monaco-collab-ext/typings/RemoteCursor"

interface Props {
	sessionID: string
	cursorManager: RemoteCursorManager
	contentManager: EditorContentManager
	event?: { NoteEvent: NoteEvent }
	setHistory(history: string[]): void
	history: string[]
}
const colours = ["#6F00E5", "#0900E6", "#005DE7", "#00C5E8", "#00E9A4", "#00EA3C", "#2BEB00", "#95EC00", "#EDDB00", "#EE7200", "#EF0700"]

const Changes = ({ sessionID, contentManager, cursorManager, history, event, setHistory }: Props) => {
	if (!event) {
		return <></>
	}

	const [cursors, setCursors] = useState<{ [sessionID: string]: RemoteCursor }>({})

	const { insert, eventID, replace, remove, cursor, userName } = event.NoteEvent
	const arrayCheck = history.indexOf(eventID)
	if (arrayCheck === -1) {
		const newHistory = [...history, eventID]
		setHistory(newHistory)

		console.debug("CollabEvent:", event)

		if (event.NoteEvent.sessionID !== sessionID) {
			if (insert) {
				contentManager.insert(insert.index, insert.text)
			}
			if (replace) {
				console.log("REPLACING!!!!")
				contentManager.replace(replace.index, replace.length, replace.text)
			}
			if (remove) {
				contentManager.delete(remove.index, remove.length)
			}
			if (cursor) {
				let editCursor: RemoteCursor
				const sID = event.NoteEvent.sessionID

				if (cursors[event.NoteEvent.sessionID]) {
					editCursor = cursors[event.NoteEvent.sessionID]
				} else {
					editCursor = cursorManager.addCursor(sID, colours[objectSize(cursors)], userName)
					let newCursors = { ...cursors }
					newCursors[sID] = editCursor
					setCursors(newCursors)
				}
				editCursor.setPosition(cursor)
			}
		}
	}

	return <></>
}

export default Changes
