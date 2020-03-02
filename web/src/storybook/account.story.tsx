// Button.story.js
import * as React from "react"
import { storiesOf } from "@storybook/react"
import { Account } from "../pages/account"

const AccountStory = () => <Account />

storiesOf("Account", module).add("default", () => <AccountStory />)
