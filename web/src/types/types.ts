export interface UserError {
	message: string
	field: string[]
}

export interface Onboard {
	prospect: Prospect
}

export interface Prospect {
	id: string
	email: string
	name: string
	onboardingComplete: boolean
}

export interface OnboardingInput {
	email: string
	name: string
}

export interface Login {
	email: string
	password: string
}

export interface User {
	id: string
	email: string
	name: string
	verified: boolean
}

export interface ErrorMap {
	[key: string]: string
}

export interface INote {
	id: string
	name: string
	body: string
	done: boolean
	owner: User
}

export interface CursorInput {
	lineNumber: number
	column: number
}

export interface CursorPlacement {
	lineNumber: number
	column: number
}

export interface NoteInsert {
	text: string
	index: number
}

export interface ReplaceText {
	text: string
	index: number
	length: number
}

export interface DeleteText {
	index: number
	length: number
}

//whenever you change this, you need to update the graphql
export interface NoteEvent {
	noteID: string
	eventID: string
	sessionID: string
	insert?: NoteInsert
	cursor?: CursorPlacement
	replace?: ReplaceText
	remove?: DeleteText
	userID: string
	userName: string
}

export interface NoteChange {
	noteID: string
	eventID: string
	sessionID: string
	insert?: NoteInsert
	replace?: ReplaceText
	remove?: DeleteText
	cursor?: CursorPlacement
}
