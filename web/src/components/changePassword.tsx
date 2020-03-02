import * as React from "react"
import { ErrorMap } from "../types/types"
import { UserContainer } from "../controllers/user"
import { ErrorBox } from "./errorBox"
import { KIND, Notification } from "baseui/notification"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { Spaced } from "./spaced"
import { Spinner } from "baseui/spinner"
import { ModalButton } from "baseui/modal"

export const ChangePassword = () => {
	const [inputError, setInputError] = React.useState<ErrorMap>({})
	const [oldPassword, setOldPassword] = React.useState<string>("")
	const [password, setPassword] = React.useState<string>("")
	const [password2, setPassword2] = React.useState<string>("")

	const { apolloError, userErrors, loading, clearErrors, useChangePassword } = UserContainer.useContainer()
	const { changePassword, changeSuccess, setChangeSuccess } = useChangePassword(oldPassword, password)

	const handleSubmit = () => {
		setChangeSuccess(false)
		clearErrors()

		const errors: ErrorMap = {}
		let foundError = false

		if (oldPassword == "") {
			errors["oldPassword"] = "Please enter your current password"
			foundError = true
		}
		if (password == "") {
			errors["password"] = "Please enter a password"
			foundError = true
		}
		if (password2 == "") {
			errors["password2"] = "Please re-enter your password"
			foundError = true
		}
		if (password !== password2) {
			errors["password2"] = "Passwords do not match"
			foundError = true
		}

		if (foundError) {
			setInputError(errors)
			return
		}

		setInputError({})
		changePassword()
	}

	return (
		<div>
			<div>
				<ErrorBox apolloError={apolloError} userErrors={userErrors} />
				{changeSuccess && <Notification kind={KIND.positive}>Your password has been updated.</Notification>}
				<FormControl label="Old Password" error={inputError["oldPassword"]}>
					<Input
						key={"oldPassword"}
						error={!!inputError["oldPassword"]}
						positive={false}
						value={oldPassword}
						type={"password"}
						onChange={e => {
							setOldPassword(e.currentTarget.value)
						}}
						placeholder={"Enter your current password"}
					/>
				</FormControl>
				<FormControl label="New Password" error={inputError["password"]}>
					<Input
						key={"password"}
						error={!!inputError["password"]}
						positive={false}
						value={password}
						type={"password"}
						onChange={e => {
							setPassword(e.currentTarget.value)
						}}
						placeholder={"Enter a new password"}
					/>
				</FormControl>
				<FormControl label="Confirm Password" error={inputError["password2"]}>
					<Input
						key={"password2"}
						error={!!inputError["password2"]}
						positive={false}
						value={password2}
						type={"password"}
						onChange={e => {
							setPassword2(e.currentTarget.value)
						}}
						placeholder={"Re-enter your new password"}
					/>
				</FormControl>
			</div>
			<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
				{loading && <Spinner />}
				<ModalButton onClick={handleSubmit}>Save</ModalButton>
			</Spaced>
		</div>
	)
}
