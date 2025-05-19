import 'package:json_annotation/json_annotation.dart';

import 'common.dart';

part 'requests.g.dart';

// "full_name": <string>,
// "birth_date": <date ISO 8601>,
// "identity": <gender>,
// "bow": <bow_class | null>,
// "rank": <sports_rank | null>,
// "region": <string | null>,
// "federation": <string | null>,
// "club": <string | null>
@JsonSerializable()
class ChangeCompetitor {
  String fullName;
  String birthDate;
  Gender identity;
  BowClass? bow;
  SportsRank? rank;
  String? region;
  String? federation;
  String? club;

  ChangeCompetitor(
    this.fullName,
    this.birthDate,
    this.identity,
    this.bow,
    this.rank,
    this.region,
    this.federation,
    this.club,
  );

  Map<String, dynamic> toJson() => _$ChangeCompetitorToJson(this);
}

// "login": <string>,
// "password": <string>
class Credentials {
  String login;
  String password;

  Credentials(this.login, this.password);
}

// "range_ordinal": <number>,
// "shots": <[ <shot> ] | null>
@JsonSerializable()
class ChangeRange {
  int rangeOrdinal;
  List<Shot>? shots;

  ChangeRange(this.rangeOrdinal, this.shots);

  Map<String, dynamic> toJson() => _$ChangeRangeToJson(this);
}

// "score": <string>,
// "priority": <bool | null>
class ChangeShootOut {
  String score;
  bool? priority;

  ChangeShootOut(this.score, this.priority);
}
