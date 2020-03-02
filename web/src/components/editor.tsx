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
import { CursorInput, Note, NoteChange, NoteEvent, NoteInsert } from "../types/types"

const SUB = gql`
	subscription onNoteEvent($noteID: String!) {
		NoteEvent(noteID: $noteID) {
			noteID
			eventID
			insert {
				text
				index
			}
			cursor {
				lineNumber
				column
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
	note: Note
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

const NoteEditor = ({ note }: Props) => {
	const editorRef = useRef<monaco.editor.ICodeEditor>()
	const [eventHistory, setHistory] = useState<{ [id: string]: boolean }>({})

	const [insertText, insData] = useMutation<{ NoteChange: boolean }, { input: NoteChange }>(NOTECHANGE)
	const { data, loading } = useSubscription<NoteEvent>(SUB, { variables: { noteID: note.id } })

	function handleEditorDidMount(_: any, editor: ICodeEditor) {
		editorRef.current = editor
		// const remoteCursorManager = new RemoteCursorManager({
		// 	editor: editor,
		// 	tooltips: true,
		// 	tooltipDuration: 2
		// })

		editor.onDidChangeCursorPosition((e: ICursorPositionChangedEvent) => {
			console.log(e)
		})

		const contentManager = new EditorContentManager({
			editor: editor as any,
			onReplace(index, length, text) {
				console.log("Replace", index, length, text)
			},
			onInsert(index, text) {
				console.log("Insert", index, text)
				const eventID = randomID()

				const newHistory = { ...eventHistory, eventID: false }
				setHistory(newHistory)

				insertText({
					variables: {
						input: {
							noteID: note.id,
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
			}
		})
		//
		// const cursor = remoteCursorManager.addCursor("jDoe", "blue", "John Doe")
		// cursor.setOffset(4)
	}

	console.log(data, loading, "subnscr")

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
			<Editor height={"20vh"} value={note.body} language="markdown" editorDidMount={handleEditorDidMount} theme={"vs-dark"} />
		</div>
	)
}

export default NoteEditor
