package sortedset

import (
	"testing"
)

func TestNonExistingKeyMemberRank(t *testing.T) {
	zset := Create()
	if _, exists := zset.GetRank("NonExistingKey", "member"); exists {
		t.Errorf("Should not fetch non existent key.")
	}
}

func TestNonAddOfDuplicateMember(t *testing.T) {
	zset := Create()
	zset.Add("zset", "member1", 0.5)
	if zset.Add("zset", "member1", 0.5) != 0 {
		t.Errorf("Should not be able to add duplicate Member member1")
	}
	if zset.Add("zset", "member1", 1.5) != 0 {
		t.Errorf("Should not be able to add duplicate Member member1")
	}
}

func TestExistingKeyMemberRank(t *testing.T) {
	zset := Create()
	members := [5]string{
		"member1",
		"member2",
		"member3",
		"member4",
		"member5",
	}
	zset.Add("zset", "member1", 0.5)
	zset.Add("zset", "member2", 1.3)
	zset.Add("zset", "member3", 2.1)
	zset.Add("zset", "member4", 2.1)
	zset.Add("zset", "member5", 4.0)

	var expectedRank uint64 = 0
	for expectedRank <= 4 {
		rank, exists := zset.GetRank("zset", members[expectedRank])
		if !exists {
			t.Errorf("Rank for member%v Should exist.", expectedRank+1)
		}
		if rank != expectedRank {
			t.Errorf("Rank should be equal to:%v but got:%v", expectedRank, rank)
		}
		expectedRank++
	}
}

func TestExistingKeyRange(t *testing.T) {
	zset := Create()
	zset.Add("zset", "member1", 0.5)
	zset.Add("zset", "member2", 1.3)
	zset.Add("zset", "member3", 2.1)
	zset.Add("zset", "member4", 2.1)
	zset.Add("zset", "member5", 4.0)

	members, scores := zset.GetMembersAndScoreInRange("zset", 0, 0)
	if len(members) != 1 {
		t.Errorf("Length of returned member should be 1 but got %v", len(members))
	}
	if len(scores) != 1 {
		t.Errorf("Length of returned scores should be 1 but got %v", len(scores))
	}

	if members[0] != "member1" && scores[0] != 0.5 {
		t.Errorf("Member should be member1 with score 0.5 but got %v, with score %v.", members[0], scores[0])
	}
	members, scores = zset.GetMembersAndScoreInRange("zset", 5, 0)
	if len(members) != 0 {
		t.Errorf("Length of returned member should be 0 but got %v", len(members))
	}
	if len(scores) != 0 {
		t.Errorf("Length of returned scores should be 0 but got %v", len(scores))
	}

	members, scores = zset.GetMembersAndScoreInRange("zset", 2, 2)
	if len(members) != 1 {
		t.Errorf("Length of returned member should be 1 but got %v", len(members))
	}
	if len(scores) != 1 {
		t.Errorf("Length of returned scores should be 1 but got %v", len(scores))
	}

	if members[0] != "member3" && scores[0] != 2.1 {
		t.Errorf("Member should be member3 with score 2.1 but got %v, with score %v.", members[0], scores[0])
	}
	members, scores = zset.GetMembersAndScoreInRange("zset", 1, 3)
	if len(members) != 3 {
		t.Errorf("Length of returned member should be 3 but got %v", len(members))
	}
	if len(scores) != 3 {
		t.Errorf("Length of returned scores should be 3 but got %v", len(scores))
	}

	if members[0] != "member2" && scores[0] != 1.3 {
		t.Errorf("Member[0] should be member2 with score 1.3 but got %v, with score %v.", members[0], scores[0])
	}
	if members[1] != "member3" && scores[1] != 2.1 {
		t.Errorf("Member[1] should be member3 with score 2.1 but got %v, with score %v.", members[0], scores[0])
	}
	if members[2] != "member4" && scores[2] != 2.1 {
		t.Errorf("Member[2] should be member4 with score 2.1 but got %v, with score %v.", members[0], scores[0])
	}

	members, scores = zset.GetMembersAndScoreInRange("zset", 4, 4)
	if len(members) != 1 {
		t.Errorf("Length of returned member should be 1 but got %v", len(members))
	}
	if len(scores) != 1 {
		t.Errorf("Length of returned scores should be 1 but got %v", len(scores))
	}

	if members[0] != "member5" && scores[0] != 4.0 {
		t.Errorf("Member should be member4 with score 4.0 but got %v, with score %v.", members[0], scores[0])
	}
	members, scores = zset.GetMembersAndScoreInRange("zset", 4, 6)
	if len(members) != 1 {
		t.Errorf("Length of returned member should be 1 but got %v", len(members))
	}
	if len(scores) != 1 {
		t.Errorf("Length of returned scores should be 1 but got %v", len(scores))
	}

	if members[0] != "member5" && scores[0] != 4.0 {
		t.Errorf("Member should be member4 with score 4.0 but got %v, with score %v.", members[0], scores[0])
	}
}

func TestExistingKeyNegativeRange(t *testing.T) {
	zset := Create()
	zset.Add("zset", "member1", 0.5)
	zset.Add("zset", "member2", 1.3)
	zset.Add("zset", "member3", 2.1)
	zset.Add("zset", "member4", 2.1)
	zset.Add("zset", "member5", 4.0)

	members, scores := zset.GetMembersAndScoreInRange("zset", -10, -5)
	if len(members) != 1 {
		t.Errorf("Length of returned member should be 1 but got %v", len(members))
	}
	if len(scores) != 1 {
		t.Errorf("Length of returned scores should be 1 but got %v", len(scores))
	}

	if members[0] != "member1" && scores[0] != 0.5 {
		t.Errorf("Member should be member1 with score 0.5 but got %v, with score %v.", members[0], scores[0])
	}
	members, scores = zset.GetMembersAndScoreInRange("zset", -2, 0)
	if len(members) != 0 {
		t.Errorf("Length of returned member should be 0 but got %v", len(members))
	}
	if len(scores) != 0 {
		t.Errorf("Length of returned scores should be 0 but got %v", len(scores))
	}

	members, scores = zset.GetMembersAndScoreInRange("zset", -3, -1)
	if len(members) != 3 {
		t.Errorf("Length of returned member should be 3 but got %v", len(members))
	}
	if len(scores) != 3 {
		t.Errorf("Length of returned scores should be 3 but got %v", len(scores))
	}

	if members[0] != "member3" && scores[0] != 2.1 {
		t.Errorf("Member[0] should be member3 with score 2.1 but got %v, with score %v.", members[0], scores[0])
	}
	if members[1] != "member4" && scores[1] != 2.1 {
		t.Errorf("Member[1] should be member4 with score 2.1 but got %v, with score %v.", members[0], scores[0])
	}
	if members[2] != "member5" && scores[2] != 4.0 {
		t.Errorf("Member[2] should be member5 with score 4.0 but got %v, with score %v.", members[0], scores[0])
	}

	members, scores = zset.GetMembersAndScoreInRange("zset", -1, -1)
	if len(members) != 1 {
		t.Errorf("Length of returned member should be 1 but got %v", len(members))
	}
	if len(scores) != 1 {
		t.Errorf("Length of returned scores should be 1 but got %v", len(scores))
	}

	if members[0] != "member5" && scores[0] != 4.0 {
		t.Errorf("Member should be member4 with score 4.0 but got %v, with score %v.", members[0], scores[0])
	}
}
