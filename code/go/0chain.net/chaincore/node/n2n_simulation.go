// +build !development

package node

//InduceDelay - induces network delay - it's a noop for production deployment
func (nd *GNode) InduceDelay(toNode *GNode) {
}

//ReadNetworkDelays - read the network delay configuration - it's a noop for production ndeployment
func ReadNetworkDelays(file string) {

}
