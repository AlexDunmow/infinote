// config.js

import { configure } from "@storybook/react"

function loadStories() {
	require("./button.story.tsx")
	require("./account.story.tsx")
	require("./login.story.tsx")
	require("./authenticate.story.tsx")
}

configure(loadStories, module)
