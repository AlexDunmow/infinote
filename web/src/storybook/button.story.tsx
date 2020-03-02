// Button.story.js

import * as React from "react"
import { storiesOf } from "@storybook/react"

const Button = ({ backgroundColor, label }: { backgroundColor: string; label: string }) => {
	return <button>{label}</button>
}

storiesOf("Button", module).add("default", () => <Button backgroundColor={"blue"} label={"Submit"} />)
