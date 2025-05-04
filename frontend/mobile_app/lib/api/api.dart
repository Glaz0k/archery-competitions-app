import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';

abstract class Api {
  /// GET /competitions/{competition_id}/competitors
  /// Получить список зарегистрированных участников.
  /// Участник имеет доступ только если сам зарегистрирован
  ///
  /// Исключения:
  /// - Соревнование не найдено
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
  /// - Невозможно добавить участников после окончания соревнования
  /// - Соревнование или участник не найдены
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
  /// - Соревнование не найдено
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitions.md#get-competitionscompetition_idindividual_groups
  Future<List<IndividualGroup>> getCompetitionsIndividualGroups(
    int competitionId,
  );

  /// POST /competitors/registration
  ///
  /// Зарегистрировать участника c помощью его аутентификационного токена
  ///
  /// Исключения:
  /// - Неверные параметры
  /// - Уже зарегистрирован
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitors.md#post-competitorsregistration
  Future<CompetitorFull> registerCompetitor(ChangeCompetitor request);

  /// GET /competitors/{competitor_id}
  ///
  /// Получить информацию об участнике
  ///
  /// Исключения:
  /// - Не существует
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitors.md#get-competitorscompetitor_id
  Future<CompetitorFull> getCompetitor(int competitorId);

  /// PUT /competitors/{competitor_id}
  ///
  /// Изменить информацию об участнике.
  /// Участнику доступно изменение только своих данных.
  /// Невозможно изменить пол и тип лука,
  /// если участник зарегистрирован в группе не на стадии создания.
  ///
  /// Исключения:
  /// - Неверные параметры
  /// - Не зарегистрирован
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitors.md#put-competitorscompetitor_id
  Future<CompetitorFull> putCompetitor(
    int competitorId,
    ChangeCompetitor request,
  );

  /// GET /cups/{cup_id}
  ///
  /// Получить кубок (участник имеет доступ, если принимает участие)
  ///
  /// Исключения:
  /// - Кубок не найден
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/cups.md#get-cupscup_id
  Future<Cup> getCup(int cupId);

  /// GET /cups
  ///
  /// Получить кубки, доступные пользователю (участник получает только те, в которых принимает участие)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/cups.md#get-cups
  Future<List<Cup>> getCups();

  /// GET /cups/{cup_id}/competitions
  ///
  /// Получить все соревнования, доступные пользователю (участник получает только те, в которых принимает участие)
  ///
  /// Исключения:
  /// - Кубок не найден
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/cups.md#get-cupscup_idcompetitions
  Future<List<Competition>> getCupsCompetitions(int cupId);
}
