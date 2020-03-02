import { ErrorBox } from "./errorBox"
import { KIND, Notification } from "baseui/notification"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import * as React from "react"
import { UserContainer } from "../controllers/user"
import { ErrorMap } from "../types/types"
import { Spinner } from "baseui/spinner"
import { ModalButton } from "baseui/modal"
import { Spaced } from "./spaced"

export const ChangeSettings = () => {
	const { apolloError, userErrors, loading, clearErrors, useChangeDetails, user } = UserContainer.useContainer()

	const [name, setName] = React.useState<string>(user ? user.name : "")
	const [inputError, setInputError] = React.useState<ErrorMap>({})
	const { changeDetails, changeSuccess, setChangeSuccess } = useChangeDetails(name)

	const handleSubmit = () => {
		setChangeSuccess(false)
		clearErrors()

		const errors: ErrorMap = {}
		let foundError = false

		if (name == "") {
			errors["name"] = "Your first name can't be empty"
			foundError = true
		}
		if (foundError) {
			setInputError(errors)
			return
		}

		setInputError({})
		changeDetails()
	}

	return (
		<div>
			<ErrorBox apolloError={apolloError} userErrors={userErrors} />
			{changeSuccess && <Notification kind={KIND.positive}>Your details have been updated.</Notification>}
			<FormControl label="Name" error={inputError["name"]}>
				<Input
					key={"name"}
					error={!!inputError["name"]}
					positive={false}
					value={name}
					onChange={e => {
						setName(e.currentTarget.value)
					}}
					placeholder={"Enter your current password"}
				/>
			</FormControl>
			<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
				{loading && <Spinner />}
				<ModalButton onClick={handleSubmit}>Save</ModalButton>
			</Spaced>
		</div>
	)
}
