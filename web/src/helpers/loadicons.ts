import { library as IconLibrary } from "@fortawesome/fontawesome-svg-core"

import {
	faHeadset,
	faHandHoldingBox,
	faUser,
	faEdit,
	faTrash,
	faMoneyCheckAlt,
	faUserCircle,
	faEnvelope,
	faChevronDown,
	faChevronUp,
	faPlus,
	faUserPlus
} from "@fortawesome/pro-solid-svg-icons"

export const loadIcons = () => {
	IconLibrary.add(
		faPlus,
		faUserPlus,
		faTrash,
		faEdit,
		faHeadset,
		faHandHoldingBox,
		faUser,
		faMoneyCheckAlt,
		faUserCircle,
		faEnvelope,
		faChevronDown,
		faChevronUp
	)
}
