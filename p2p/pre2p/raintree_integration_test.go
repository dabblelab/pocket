package pre2p

import (
	"sync"
	"testing"
	"time"

	typesPre2P "github.com/pokt-network/pocket/p2p/pre2p/types"
	"github.com/pokt-network/pocket/shared/modules"
	"github.com/pokt-network/pocket/shared/types"
	"google.golang.org/protobuf/types/known/anypb"
)

// IMPROVE(team): Looking into adding more tests and accounting for more edge cases.

func TestRainTreeCompleteOneNodes(t *testing.T) {
	// val_1
	originatorNode := validatorId(t, 1)
	var expectedCalls = TestRainTreeCommConfig{
		validatorId(t, 1): {0, 0}, // {numReads, numWrites}
	}
	testRainTreeCalls(t, originatorNode, expectedCalls, false)
}

func TestRainTreeCompleteTwoNodes(t *testing.T) {
	// val_1
	//   └───────┐
	// 	       val_2
	originatorNode := validatorId(t, 1)
	var expectedCalls = TestRainTreeCommConfig{
		validatorId(t, 1): {0, 0}, // Originator
		validatorId(t, 2): {1, 1},
	}
	testRainTreeCalls(t, originatorNode, expectedCalls, false)
}

func TestRainTreeCompleteThreeNodes(t *testing.T) {
	// 	          val_1
	// 	   ┌───────┴────┬─────────┐
	//   val_2        val_1     val_3
	originatorNode := validatorId(t, 1)
	var expectedCalls = TestRainTreeCommConfig{
		validatorId(t, 1): {0, 0}, // Originator
		validatorId(t, 2): {1, 1},
		validatorId(t, 3): {1, 1},
	}
	testRainTreeCalls(t, originatorNode, expectedCalls, false)
}

func TestRainTreeCompleteFourNodes(t *testing.T) {
	// Test configurations (visualization retrieved from simulator)
	// 	                val_1
	// 	  ┌───────────────┴────┬─────────────────┐
	//  val_2                val_1             val_3
	//    └───────┐            └───────┐         └───────┐
	// 		    val_3                val_2             val_4
	originatorNode := validatorId(t, 1)
	var expectedCalls = TestRainTreeCommConfig{
		validatorId(t, 1): {0, 0}, // Originator
		validatorId(t, 2): {2, 2},
		validatorId(t, 3): {2, 2},
		validatorId(t, 4): {1, 1},
	}
	testRainTreeCalls(t, originatorNode, expectedCalls, false)
}

func TestRainTreeCompleteNineNodes(t *testing.T) {
	// 	                              val_1
	// 	         ┌──────────────────────┴────────────┬────────────────────────────────┐
	//         val_4                               val_1                            val_7
	//   ┌───────┴────┬─────────┐            ┌───────┴────┬─────────┐         ┌───────┴────┬─────────┐
	// val_6        val_4     val_8        val_3        val_1     val_5     val_9        val_7     val_2
	originatorNode := validatorId(t, 1)
	var expectedCalls = TestRainTreeCommConfig{
		validatorId(t, 1): {0, 0}, // Originator
		validatorId(t, 2): {1, 1},
		validatorId(t, 3): {1, 1},
		validatorId(t, 4): {1, 1},
		validatorId(t, 5): {1, 1},
		validatorId(t, 6): {1, 1},
		validatorId(t, 7): {1, 1},
		validatorId(t, 8): {1, 1},
		validatorId(t, 9): {1, 1},
	}
	testRainTreeCalls(t, originatorNode, expectedCalls, false)
}

func TestRainTreeCompleteEighteenNodes(t *testing.T) {
	// 	                                                                                                              val_1
	// 	                                      ┌──────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────────┐
	//                                      val_7                                                                                                            val_1                                                                                                     val_13
	//             ┌──────────────────────────┴────────────┬────────────────────────────────────┐                                     ┌────────────────────────┴────────────┬──────────────────────────────────┐                                ┌────────────────────────┴──────────────┬────────────────────────────────────┐
	//           val_11                                   val_7                               val_15                                 val_5                                 val_1                              val_9                           val_17                                  val_13                                val_3
	//    ┌────────┴─────┬───────────┐             ┌───────┴────┬──────────┐           ┌────────┴─────┬──────────┐            ┌───────┴────┬──────────┐             ┌───────┴────┬─────────┐          ┌────────┴────┬─────────┐         ┌───────┴─────┬──────────┐             ┌────────┴─────┬───────────┐          ┌───────┴────┬──────────┐
	// val_13         val_11      val_16        val_9        val_7      val_12      val_17         val_15     val_8        val_7        val_5      val_10        val_3        val_1     val_6      val_11        val_9     val_2     val_1         val_17     val_4         val_15         val_13      val_18     val_5        val_3      val_14
	originatorNode := validatorId(t, 1)
	var expectedCalls = TestRainTreeCommConfig{
		validatorId(t, 1):  {1, 1}, // Originator
		validatorId(t, 2):  {1, 1},
		validatorId(t, 3):  {2, 2},
		validatorId(t, 4):  {1, 1},
		validatorId(t, 5):  {2, 2},
		validatorId(t, 6):  {1, 1},
		validatorId(t, 7):  {2, 2},
		validatorId(t, 8):  {1, 1},
		validatorId(t, 9):  {2, 2},
		validatorId(t, 10): {1, 1},
		validatorId(t, 11): {2, 2},
		validatorId(t, 12): {1, 1},
		validatorId(t, 13): {2, 2},
		validatorId(t, 14): {1, 1},
		validatorId(t, 15): {2, 2},
		validatorId(t, 16): {1, 1},
		validatorId(t, 17): {2, 2},
		validatorId(t, 18): {1, 1},
	}
	// Note that the originator, `val_1` is also messaged by `val_17` outside of continuously
	// demoting itself.
	testRainTreeCalls(t, originatorNode, expectedCalls, true)
}

func TestRainTreeCompleteTwentySevenNodes(t *testing.T) {
	// 	                                                                                                                    val_1
	// 	                                     ┌────────────────────────────────────────────────────────────────────────────────┴───────────────────────────────────────┬───────────────────────────────────────────────────────────────────────────────────────────────────────────┐
	//                                    val_10                                                                                                                   val_1                                                                                                       val_19
	//            ┌──────────────────────────┴──────────────┬──────────────────────────────────────┐                                         ┌────────────────────────┴────────────┬──────────────────────────────────┐                                  ┌────────────────────────┴──────────────┬────────────────────────────────────┐
	//          val_16                                    val_10                                 val_22                                     val_7                                 val_1                             val_13                             val_25                                  val_19                                val_4
	//   ┌────────┴─────┬───────────┐              ┌────────┴─────┬───────────┐           ┌────────┴─────┬───────────┐              ┌────────┴────┬──────────┐             ┌───────┴────┬─────────┐          ┌────────┴─────┬──────────┐         ┌───────┴─────┬──────────┐             ┌────────┴─────┬───────────┐          ┌───────┴────┬──────────┐
	// val_20         val_16      val_24         val_14         val_10      val_18      val_26         val_22      val_12         val_11        val_7      val_15        val_5        val_1     val_9      val_17         val_13     val_3     val_2         val_25     val_6         val_23         val_19      val_27     val_8        val_4      val_21
	originatorNode := validatorId(t, 1)
	var expectedCalls = TestRainTreeCommConfig{
		validatorId(t, 1):  {0, 0}, // Originator
		validatorId(t, 2):  {1, 1},
		validatorId(t, 3):  {1, 1},
		validatorId(t, 4):  {1, 1},
		validatorId(t, 5):  {1, 1},
		validatorId(t, 6):  {1, 1},
		validatorId(t, 7):  {1, 1},
		validatorId(t, 8):  {1, 1},
		validatorId(t, 9):  {1, 1},
		validatorId(t, 10): {1, 1},
		validatorId(t, 11): {1, 1},
		validatorId(t, 12): {1, 1},
		validatorId(t, 13): {1, 1},
		validatorId(t, 14): {1, 1},
		validatorId(t, 15): {1, 1},
		validatorId(t, 16): {1, 1},
		validatorId(t, 17): {1, 1},
		validatorId(t, 18): {1, 1},
		validatorId(t, 19): {1, 1},
		validatorId(t, 20): {1, 1},
		validatorId(t, 21): {1, 1},
		validatorId(t, 22): {1, 1},
		validatorId(t, 23): {1, 1},
		validatorId(t, 24): {1, 1},
		validatorId(t, 25): {1, 1},
		validatorId(t, 26): {1, 1},
		validatorId(t, 27): {1, 1},
	}
	testRainTreeCalls(t, originatorNode, expectedCalls, false)
}

func testRainTreeCalls(t *testing.T, origNode string, testCommConfig TestRainTreeCommConfig, isOriginatorPinged bool) {
	// Network configurations
	numValidators := len(testCommConfig)
	configs := createConfigs(t, numValidators)

	// Test configurations
	var messageHandeledWaitGroup sync.WaitGroup
	if isOriginatorPinged {
		messageHandeledWaitGroup.Add(numValidators)
	} else {
		messageHandeledWaitGroup.Add(numValidators - 1) // -1 because the originator node implicitly handles the message
	}

	// Network initialization
	connMocks := make(map[string]typesPre2P.TransportLayerConn)
	busMocks := make(map[string]modules.Bus)
	for valId, expectedCall := range testCommConfig {
		connMocks[valId] = prepareConnMock(t, expectedCall.numNetworkReads, expectedCall.numNetworkWrites)
		busMocks[valId] = prepareBusMock(t, &messageHandeledWaitGroup)
	}

	// Module injection
	p2pModules := prepareP2PModules(t, configs)
	for validatorId, mod := range p2pModules {
		mod.listener = connMocks[validatorId]
		mod.SetBus(busMocks[validatorId])
		for _, peer := range mod.network.GetAddrBook() {
			peer.Dialer = connMocks[peer.ServiceUrl]
		}
		mod.Start()
		defer mod.Stop()
	}

	// Trigger originator message
	p := &anypb.Any{}
	p2pMod := p2pModules[origNode]
	p2pMod.Broadcast(p, types.PocketTopic_DEBUG_TOPIC)

	// Wait for completion
	done := make(chan struct{})
	go func() {
		messageHandeledWaitGroup.Wait()
		close(done)
	}()

	// Timeout or succeed
	select {
	case <-done:
	// All done!
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for message to be handled")
	}
}
