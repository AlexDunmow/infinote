import * as React from "react"
import { AuthContainer } from "../controllers/auth"
import { Modal, ModalBody, ModalHeader, ModalFooter, ModalButton, ROLE, SIZE } from "baseui/modal"
import { ErrorBox } from "../components/errorBox"
import { FormControl } from "baseui/form-control"
import { Input } from "baseui/input"
import { UserContainer } from "../controllers/user"
import { Spinner } from "baseui/spinner"
import { Spaced } from "../components/spaced"
import { ErrorMap } from "../types/types"

export const VerificationComplete = () => {
	const { setShowVerifyComplete, showVerifyComplete, apolloError, userErrors, loading } = AuthContainer.useContainer()

	const [inputError, setInputError] = React.useState<ErrorMap>({})
	const [password, setPassword] = React.useState<string>("")
	const [password2, setPassword2] = React.useState<string>("")

	const { useChangePassword } = UserContainer.useContainer()
	const { changePassword, changeSuccess } = useChangePassword("", password)

	const handleSubmit = () => {
		const errors: ErrorMap = {}
		let foundError = false

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
		<Modal closeable={false} isOpen={showVerifyComplete} animate size={SIZE.default} role={ROLE.dialog}>
			<ModalHeader>{!changeSuccess ? "Verification Complete" : "Success!"}</ModalHeader>

			{!changeSuccess ? (
				<ModalBody>
					Thanks for verifying your email address. Before you go further, let's set up a password.
					<div>
						<ErrorBox apolloError={apolloError} userErrors={userErrors} />
						<FormControl label="Password" error={inputError["password"]}>
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
				</ModalBody>
			) : (
				<ModalBody>Your password has been set. Click the close button below to continue.</ModalBody>
			)}
			<ModalFooter>
				{!changeSuccess ? (
					<Spaced overrides={{ container: { justifyContent: "flex-end" } }}>
						{loading && <Spinner />}
						<ModalButton onClick={handleSubmit}>Save</ModalButton>
					</Spaced>
				) : (
					<ModalButton onClick={() => setShowVerifyComplete(false)}>Close</ModalButton>
				)}
			</ModalFooter>
		</Modal>
	)
}
