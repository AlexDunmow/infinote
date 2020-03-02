import { library as IconLibrary } from "@fortawesome/fontawesome-svg-core"

import { faHeadset, faHandHoldingBox, faUser, faMoneyCheckAlt, faUserCircle, faEnvelope, faChevronDown, faChevronUp } from "@fortawesome/pro-light-svg-icons"

export const loadIcons = () => {
	IconLibrary.add(faHeadset, faHandHoldingBox, faUser, faMoneyCheckAlt, faUserCircle, faEnvelope, faChevronDown, faChevronUp)
}
