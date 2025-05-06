import 'common.dart';

// "id": <number>,
// "full_name": <string>,
// "birth_date": <date ISO 8601>,
// "identity": <gender>,
// "bow": <bow_class | null>,
// "rank": <sports_rank | null>,
// "region": <string | null>,
// "federation": <string | null>,
// "club": <string | null>
class CompetitorFull {
  final int id;
  final String fullName;
  final DateTime birthDate;
  final Gender identity;
  final BowClass? bow;
  final SportsRank? rank;
  final String? region;
  final String? federation;
  final String? club;

  CompetitorFull(
    this.id,
    this.fullName,
    this.birthDate,
    this.identity,
    this.bow,
    this.rank,
    this.region,
    this.federation,
    this.club,
  );
}

// "competition_id": <number>,
// "competitor": <competitor_full>,
// "is_active": <bool>,
// "created_at": <YYYY-MM-DDThh:mm:ss±hh ISO 8601>
class CompetitorCompetitionDetail {
  final int competitionId;
  final CompetitorFull competitor;
  final bool isActive;
  final DateTime createdAt;

  CompetitorCompetitionDetail(
    this.competitionId,
    this.competitor,
    this.isActive,
    this.createdAt,
  );
}

// "id": <number>,
// "competition_id": <number>,
// "bow": <bow_class>,
// "identity": <gender | null>,
// "state": <group_state>
class IndividualGroup {
  final int id;
  final int competitionId;
  final BowClass bow;
  final Gender? identity;
  final GroupState state;

  IndividualGroup(
    this.id,
    this.competitionId,
    this.bow,
    this.identity,
    this.state,
  );
}

// "id": <number>,
// "title": <string>,
// "address": <string | null>,
// "season": <string | null>
class Cup {
  final int id;
  final String title;
  final String? address;
  final String? season;

  Cup(this.id, this.title, this.address, this.season);
}

// "id": <number>,
// "stage": <competition_stage>,
// "start_date": <YYYY-MM-DD ISO 8601 | null>,
// "end_date": <YYYY-MM-DD ISO 8601 | null>,
// "is_ended": <bool>
class Competition {
  final int id;
  final CompetitionStage stage;
  final DateTime? startDate;
  final DateTime? endDate;
  final bool isEnded;

  Competition(this.id, this.stage, this.startDate, this.endDate, this.isEnded);
}

// "group_id": <number>,
// "competitor": <competitor_full>
class CompetitorGroupDetail {
  final int groupId;
  final CompetitorFull competitor;

  CompetitorGroupDetail(this.groupId, this.competitor);
}

// "group_id": <number>,
// "distance": <string>,
// "round_count": <number>,
// "sections": [ <section> ]
class QualificationTable {
  final int groupId;
  final String distance;
  final int roundCount;
  final List<Section> sections;

  QualificationTable(
    this.groupId,
    this.distance,
    this.roundCount,
    this.sections,
  );
}

// "id": <number>,
// "competitor": <competitor_shrinked>,
// "place": <number | null>,
// "rounds": [ <qualification_round_shrinked> ],
// "total": <number | null>,
// "10_s": <number | null>,
// "9_s": <number | null>,
// "rank_gained": <sports_rank | null>
class Section {
  final int id;
  final CompetitorShrinked competitor;
  final int? place;
  final List<QualificationRoundShrinked> rounds;
  final int? total;
  final int? tenS;
  final int? nineS;
  final SportsRank? rankGained;

  Section(
    this.id,
    this.competitor,
    this.place,
    this.rounds,
    this.total,
    this.tenS,
    this.nineS,
    this.rankGained,
  );
}

// "id": <number>,
// "full_name": <string>
class CompetitorShrinked {
  final int id;
  final String fullName;

  CompetitorShrinked(this.id, this.fullName);
}

// "round_ordinal": <number>,
// "is_active": <bool>,
// "total_score": <number | null>
class QualificationRoundShrinked {
  final int roundOrdinal;
  final bool isActive;
  final int? totalScore;

  QualificationRoundShrinked(this.roundOrdinal, this.isActive, this.totalScore);
}

// "group_id": <number>,
// "quarterfinal": <quarterfinal>,
// "semifinal": <semifinal | null>,
// "final": <final | null>
class FinalGrid {
  final int groupId;
  final Quarterfinal quarterfinal;
  final Semifinal? semifinal;
  final Final? fina1;

  FinalGrid(this.groupId, this.quarterfinal, this.semifinal, this.fina1);
}

// "sparring_1": <sparring | null>,
// "sparring_2": <sparring | null>,
// "sparring_3": <sparring | null>,
// "sparring_4": <sparring | null>
class Quarterfinal {
  final Sparring? sparring1;
  final Sparring? sparring2;
  final Sparring? sparring3;
  final Sparring? sparring4;

  Quarterfinal(this.sparring1, this.sparring2, this.sparring3, this.sparring4);
}

// "sparring_5": <sparring | null>,
// "sparring_6": <sparring | null>
class Semifinal {
  final Sparring? sparring5;
  final Sparring? sparring6;

  Semifinal(this.sparring5, this.sparring6);
}

// "sparring_gold": <sparring | null>,
// "sparring_bronze": <sparring | null>
class Final {
  final Sparring? sparringGold;
  final Sparring? sparringBronze;

  Final(this.sparringGold, this.sparringBronze);
}

// "id": <number>,
// "top_place": <sparring_place | null>,
// "bot_place": <sparring_place | null>,
// "state": <sparring_state>
class Sparring {
  final int number;
  final SparringPlace? topPlace;
  final SparringPlace? botPlace;
  final SparringState state;

  Sparring(this.number, this.topPlace, this.botPlace, this.state);
}

// "id": <number>,
// "competitor": <competitor_shrinked>,
// "range_group": <range_group>,
// "is_active": <bool>,
// "shoot_out": <shoot_out | null>,
// "sparring_score": <number>
class SparringPlace {
  final int number;
  final CompetitorShrinked competitor;
  final RangeGroup rangeGroup;
  final bool isActive;
  final ShootOut? shootOut;
  final int sparringScore;

  SparringPlace(
    this.number,
    this.competitor,
    this.rangeGroup,
    this.isActive,
    this.shootOut,
    this.sparringScore,
  );
}

// "id": <number>,
// "ranges_max_count": <number>,
// "range_size": <number>,
// "type": <range_type>
// "ranges": <[ <range> ]>,
// "total_score": <number | null>
class RangeGroup {
  final int number;
  final int rangesMaxCount;
  final int rangeSize;
  final RangeType type;
  final List<Range> ranges;
  final int? totalScore;

  RangeGroup(
    this.number,
    this.rangesMaxCount,
    this.rangeSize,
    this.type,
    this.ranges,
    this.totalScore,
  );
}

// "id": <number>,
// "range_ordinal": <number>,
// "is_active": <bool>,
// "shots": <[ <shot> ] | null>,
// "range_score": <number | null>
class Range {
  final int id;
  final int rangeOrdinal;
  final bool isActive;
  final List<Shot>? shots;
  final int? rangeScore;

  Range(this.id, this.rangeOrdinal, this.isActive, this.shots, this.rangeScore);
}

// id": <number>,
// "score": <string>,
// "priority": <bool | null>
class ShootOut {
  final int number;
  final String score;
  final bool? priority;

  ShootOut(this.number, this.score, this.priority);
}

// "section_id": <number>,
// "round_ordinal": <number>,
// "is_active": <bool>,
// "range_group": <range_group>
class QualificationRoundFull {
  final int sectionId;
  final int roundOrdinal;
  final bool isActive;
  final RangeGroup rangeGroup;

  QualificationRoundFull(
    this.sectionId,
    this.roundOrdinal,
    this.isActive,
    this.rangeGroup,
  );
}

// "id": <number>,
// "competitor": <competitor_shrinked>,
// "range_group": <range_group>,
// "is_active": <bool>,
// "shoot_out": <shoot_out | null>,
// "sparring_score": <number>
class SparingPlace {
  final int id;
  final CompetitorShrinked competitor;
  final RangeGroup rangeGroup;
  final bool isActive;
  final ShootOut? shootOut;
  final int sparringScore;

  SparingPlace(
    this.id,
    this.competitor,
    this.rangeGroup,
    this.isActive,
    this.shootOut,
    this.sparringScore,
  );
}
