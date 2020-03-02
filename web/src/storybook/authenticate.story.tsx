// Button.story.js
import * as React from "react"
import { storiesOf } from "@storybook/react"
import { Provider as StyletronProvider } from "styletron-react"
import { Authenticate } from "../components/authenticate"
import { MemoryRouter } from "react-router"
import { AuthContainer } from "../controllers/auth"
import { BaseProvider, LightTheme } from "baseui"
import { Client as Styletron } from "styletron-engine-atomic"
const engine = new Styletron()

storiesOf("Authenticate", module)
	.addDecorator(story => <MemoryRouter initialEntries={["/"]}>{story()}</MemoryRouter>)
	// .addDecorator((story) =><AuthContainer.Provider>{story()}</AuthContainer.Provider>)
	.addDecorator(story => <StyletronProvider value={engine}>{story()}</StyletronProvider>)
	.addDecorator(story => <BaseProvider theme={LightTheme}>{story()}</BaseProvider>)
	.add("default", () => <Authenticate />)
