import * as React from "react"
import * as monaco from "monaco-editor"
import Editor, { Monaco } from "@monaco-editor/react"
import { useRef, useState } from "react"
import * as monacoEditor from "monaco-editor"
import { EditorContentManager, RemoteCursorManager } from "@convergencelabs/monaco-collab-ext"
import { gql } from "apollo-boost"
import { editor } from "monaco-editor"
import ICursorPositionChangedEvent = editor.ICursorPositionChangedEvent
import ICodeEditor = editor.ICodeEditor
import { useMutation, useSubscription } from "@apollo/react-hooks"
import { CursorInput, INote, NoteChange, NoteEvent, NoteInsert } from "../../types/types"
import Changes from "./changes"

const SUB = gql`
	subscription onNoteEvent($noteID: String!, $sessionID: String!) {
		NoteEvent(noteID: $noteID, sessionID: $sessionID) {
			noteID
			eventID
			sessionID
			insert {
				text
				index
			}
			cursor {
				lineNumber
				column
			}
			replace {
				length
				index
				text
			}
			remove {
				length
				index
			}
			userID
			userName
		}
	}
`

const NOTECHANGE = gql`
	mutation NoteChange($input: NoteChange!) {
		NoteChange(input: $input) {
			success
		}
	}
`

interface Props {
	note: INote
}

function randomID(): string {
	// Math.random should be unique because of its seeding algorithm.
	// Convert it to base 36 (numbers + letters), and grab the first 9 characters
	// after the decimal.
	return (
		"_" +
		Math.random()
			.toString(36)
			.substr(2, 9)
	)
}

const sessionID = randomID()

const NoteEditor = ({ note }: Props) => {
	const editorRef = useRef<monaco.editor.ICodeEditor>()
	const [history, setHistory] = useState<string[]>([])

	const [contentManager, setContentManager] = useState<EditorContentManager>()
	const [cursorManager, setCursorManager] = useState<RemoteCursorManager>()

	const [changeNote, insData] = useMutation<{ NoteChange: boolean }, { input: NoteChange }>(NOTECHANGE)
	const { data, loading } = useSubscription<{ NoteEvent: NoteEvent }>(SUB, { variables: { noteID: note.id, sessionID } })

	function handleEditorDidMount(_: any, editor: ICodeEditor) {
		editorRef.current = editor
		// const remoteCursorManager = new RemoteCursorManager({
		// 	editor: editor,
		// 	tooltips: true,
		// 	tooltipDuration: 2
		// })

		editor.onDidChangeCursorPosition((e: ICursorPositionChangedEvent) => {
			console.log("Position change", e)
			const eventID = randomID()
			changeNote({
				variables: {
					input: {
						noteID: note.id,
						sessionID,
						eventID,
						cursor: e.position
					}
				}
			})
		})

		const cManager = new EditorContentManager({
			editor: editor as any,
			onReplace(index, length, text) {
				console.log("Replace", index, length, text)
				const eventID = randomID()
				changeNote({
					variables: {
						input: {
							noteID: note.id,
							sessionID,
							eventID,
							replace: {
								text,
								length,
								index
							}
						}
					}
				})
			},
			onInsert(index, text) {
				console.log("Insert", index, text)
				const eventID = randomID()

				changeNote({
					variables: {
						input: {
							noteID: note.id,
							sessionID,
							eventID,
							insert: {
								text,
								index
							}
						}
					}
				})
			},
			onDelete(index, length) {
				console.log("Delete", index, length)
				const eventID = randomID()
				changeNote({
					variables: {
						input: {
							noteID: note.id,
							sessionID,
							eventID,
							remove: {
								length,
								index
							}
						}
					}
				})
			}
		})

		const cursManager = new RemoteCursorManager({
			editor: editor as any,
			tooltips: true,
			tooltipDuration: 2
		})

		setCursorManager(cursManager)
		setContentManager(cManager)
		//
		// const cursor = remoteCursorManager.addCursor("jDoe", "blue", "John Doe")
		// cursor.setOffset(4)
	}

	function listenEditorChanges() {
		if (editorRef.current) {
			editorRef.current.onDidChangeModelContent((ev: any) => {
				if (editorRef.current) {
					console.log(editorRef.current.getValue())
				}
			})
		}
	}
	return (
		<div>
			{contentManager && cursorManager && (
				<Changes cursorManager={cursorManager} contentManager={contentManager} event={data} setHistory={setHistory} history={history} sessionID={sessionID} />
			)}
			<Editor height={"20vh"} value={note.body} language="markdown" editorDidMount={handleEditorDidMount} theme={"vs-dark"} />
		</div>
	)
}

export default NoteEditor
