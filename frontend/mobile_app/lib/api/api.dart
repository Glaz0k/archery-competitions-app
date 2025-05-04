import 'package:mobile_app/api/data.dart';

abstract class Api {
  /// GET /competitions/{competition_id}/competitors
  /// Получить список зарегистрированных участников.
  /// Участник имеет доступ только если сам зарегистрирован
  ///
  /// Исключения:
  ///   Соревнование не найдено
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitions.md#get-competitionscompetition_idcompetitors
  Future<List<CompetitorCompetitionDetail>> getCompetitionsCompetitors(
    int competitionId,
  );

  /// PUT /competitions/{competition_id}/competitors/{competitor_id}
  ///
  /// Изменить статус участника в соревновании. Невозможно поменять статус,
  /// если участник состоит в активной группе этого
  /// соревнования(не на стадии создания).
  /// Участник имеет доступ, если меняет свой статус
  ///
  /// Исключения:
  ///   Невозможно добавить участников после окончания соревнования
  ///   Соревнование или участник не найдены
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitions.md#put-competitionscompetition_idcompetitorscompetitor_id
  Future<CompetitorCompetitionDetail> changeCompetitorStatus(
    int competitionId,
    int competitorId,
    bool status,
  );

  /// GET /competitions/{competition_id}/individual_groups
  ///
  /// Получить все группы, доступные пользователю
  /// (участник получает только те, в которых принимает участие)
  ///
  /// Исключения:
  ///   Соревнование не найдено
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitions.md#get-competitionscompetition_idindividual_groups
  Future<List<IndividualGroup>> getCompetitionsIndividualGroups(int competitionId);
}
