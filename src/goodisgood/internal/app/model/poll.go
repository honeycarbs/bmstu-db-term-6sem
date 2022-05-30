package model

type Answer struct {
	Word string
	Mark int
}

type Stats struct {
	Word  string
	Stats float64
}

type Poll struct {
	Answer []Answer
}

// CALL apoc.trigger.add(
//     'deleteAccountWithUser',
//     'UNWIND $deletedNodes as acc
//     match (u:user)-[o:OWNS]->(:acc)
//     detach delete u',
//     {phase:'before'}
// )
