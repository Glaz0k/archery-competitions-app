import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/api/responses.dart';


abstract class Api {
  /// POST /auth/registration
  ///
  /// Зарегистрировать пользователя в системе. Выдаётся внутрення роль user.
  /// После регистрации на клиент возвращается сессионная cookie
  ///
  /// Исключения:
  /// - InvalidParametersException(Неверные параметры регистрации)
  /// - AlreadyExistException(Пользователь с таким login уже зарегистрирован)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/auth.md#post-authregistration
  Future<void> register(Credentials credentials);

  /// POST /auth/login
  ///
  /// Войти с учётной записью пользователя.
  /// После входа на клиент возвращается сессионная cookie
  ///
  /// Исключения:
  /// - InvalidParametersException(Неверные параметры входа)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/auth.md#post-authlogin
  Future<void> login(Credentials credentials);

  /// POST /auth/logout
  ///
  /// Выйти из учётной записи пользователя.
  /// После выхода удаляется связанный с cookie сессионный токен и
  /// для дальнейшей работы необходимо заново войти
  ///
  /// Исключения:
  /// - NotFoundException(Пользовательская сессия не найдена)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/auth.md#post-authlogout
  Future<void> logout();

  /// GET /competitions/{competition_id}/competitors
  /// Получить список зарегистрированных участников.
  /// Участник имеет доступ только если сам зарегистрирован
  ///
  /// Исключения:
  /// - NotFoundException(Соревнование не найдено)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitions.md#get-competitionscompetition_idcompetitors
  Future<List<CompetitorCompetitionDetail>> getCompetitionsCompetitors(int competitionId);

  /// PUT /competitions/{competition_id}/competitors/{competitor_id}
  ///
  /// Изменить статус участника в соревновании. Невозможно поменять статус,
  /// если участник состоит в активной группе этого
  /// соревнования(не на стадии создания).
  /// Участник имеет доступ, если меняет свой статус
  ///
  /// Исключения:
  /// - BadActionException(Невозможно добавить участников после окончания соревнования)
  /// - NotFoundException(Соревнование или участник не найдены)
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
  /// - NotFoundException(Соревнование не найдено)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitions.md#get-competitionscompetition_idindividual_groups
  Future<List<IndividualGroup>> getCompetitionsIndividualGroups(int competitionId);

  /// POST /competitors/registration
  ///
  /// Зарегистрировать участника c помощью его аутентификационного токена
  ///
  /// Исключения:
  /// - InvalidParametersException(Неверные параметры)
  /// - AlreadyExistException(Уже зарегистрирован)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/competitors.md#post-competitorsregistration
  Future<CompetitorFull> registerCompetitor(ChangeCompetitor request);

  /// GET /competitors/{competitor_id}
  ///
  /// Получить информацию об участнике
  ///
  /// Исключения:
  /// - NotFoundException(Не существует)
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
  /// - InvalidParametersException(Неверные параметры)
  /// - NotFoundException(Не зарегистрирован)
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
  /// - NotFoundException(Кубок не найден)
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
  /// - NotFoundException(Кубок не найден)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/cups.md#get-cupscup_idcompetitions
  Future<List<Competition>> getCupsCompetitions(int cupId);

  /// GET /individual_groups/{group_id}
  ///
  /// Получить информацию о группе. Участникам доступны только те, в которых они участвуют
  ///
  /// Исключения:
  /// - NotFoundException(Группа не найдена)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/individual_groups.md#get-individual_groupsgroup_id
  Future<IndividualGroup> getIndividualGroup(int groupId);

  /// GET /individual_groups/{group_id}/competitors
  ///
  /// Получить список участников в группе. Участник имеет доступ только если состоит в группе
  ///
  /// Исключения:
  /// - NotFoundException(Группа не найдена)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/individual_groups.md#get-individual_groupsgroup_idcompetitors
  Future<List<CompetitorGroupDetail>> getIndividualGroupCompetitors(
    int groupId,
  );

  /// GET /individual_groups/{group_id}/qualification
  ///
  /// Получить текущую таблицу квалификаций для группы. Участник имеет доступ только если принадлежит к группе
  ///
  /// Исключения:
  /// - NotFoundException(Группа или квалификация не найдена)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/individual_groups.md#get-individual_groupsgroup_idqualification
  Future<QualificationTable> getIndividualGroupQualificationTable(int groupId);

  /// GET /individual_groups/{group_id}/final_grid
  ///
  /// Получить финальную сетку группы. Участник имеет доступ только если принадлежит к группе
  ///
  /// Исключения:
  /// - NotFoundException(Группа или сетка не найдена)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/individual_groups.md#get-individual_groupsgroup_idfinal_grid
  Future<FinalGrid> getIndividualGroupFinalGrid(int groupId);

  /// GET /qualification_sections/{id}
  ///
  /// Получить секцию. Участник имеет доступ только если это его секция
  ///
  /// Исключения:
  /// - NotFoundException(Секция не найдена)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/qualification_sections.md#get-qualification_sectionsid
  Future<Section> getQualificationSection(int sectionId);

  /// GET /qualification_sections/{id}/rounds/{round_ordinal}
  ///
  /// Получить раунд. Участник имеет доступ только если это его секция
  ///
  /// Исключения:
  /// - NotFoundException(Секция или раунд не найдены)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/qualification_sections.md#get-qualification_sectionsidroundsround_ordinal
  Future<QualificationRoundFull> getQualificationSectionsRound(
    int sectionId,
    int roundOrdinal,
  );

  /// GET /qualification_sections/{id}/rounds/{round_ordinal}/ranges
  ///
  /// Получить серии. Участник имеет доступ только если это его секция
  ///
  /// Исключения:
  ///  - NotFoundException(Секция или раунд не найдены)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/qualification_sections.md#get-qualification_sectionsidroundsround_ordinalranges
  Future<RangeGroup> getQualificationSectionsRoundsRanges(
    int sectionId,
    int roundOrdinal,
  );

  /// PUT /qualification_sections/{id}/rounds/{round_ordinal}/ranges
  ///
  /// Изменить серию. Участник имеет доступ только если это его секция и серия активна. Если номера выстрелов в списке дублируются применяется последний встретившийся по порядку. При изменении организатором завершённой серии, результаты должны повлиять на итоговый счёт.
  ///
  /// Important
  ///
  /// Изменение серии - транзакция, если во время изменения произошла ошибка, частичные изменения недопустимы
  ///
  /// Исключения:
  /// - NotFoundException(Секция, раунд или серия не найдены)
  /// - InvalidParametersException(Неверные параметры)
  /// - InvalidScoreException(Счет выстрела не соответствует типу серии)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/qualification_sections.md#put-qualification_sectionsidroundsround_ordinalranges
  Future<Range> putQualificationSectionsRoundsRange(
    int sectionId,
    int roundOrdinal,
    ChangeRange request,
  );

  /// POST /qualification_sections/{id}/rounds/{round_ordinal}/ranges/{range_ordinal}/end
  ///
  /// Подтвердить, что серия закончена.
  ///
  /// Если серия не является последней в раунде, необходимо привести в активное состояние следующую.
  /// Иначе, необходимо завершить раунд.
  /// Если завершенный раунд не последний, привести в активное состояние следующий раунд и его первую серию.
  ///
  /// Участник имеет доступ только если это его секция.
  /// Серию можно завершить, если у всех выстрелов в серии выставлен non-null счёт.
  /// Завершение неактивной серии невозможно
  ///
  /// Исключения:
  /// - NotFoundException(Секция, раунд или серия не найдены)
  /// - BadActionException(Невозможно завершить неполную или неактивную серию)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/qualification_sections.md#post-qualification_sectionsidroundsround_ordinalrangesrange_ordinalend
  Future<Range> endQualificationSectionsRoundsRange(
    int sectionId,
    int roundOrdinal,
    int rangeOrdinal,
  );

  /// GET /sparring_places/{id}
  ///
  /// Получить место. Участник имеет доступ только если это его место. При вычислении счёта по типу "победные очки" необходимо учитывать только серии, которые оба участника завершили
  ///
  /// Исключения:
  /// - NotFoundException(Место не найдено)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/sparring_places.md#get-sparring_placesid
  Future<SparingPlace> getSparringPlace(int placeId);

  /// GET /sparring_places/{id}/ranges
  ///
  /// Получить серии. Участник имеет доступ только если это его место
  ///
  /// Исключения:
  /// - NotFoundException(Место или серия не найдены)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/sparring_places.md#get-sparring_placesidranges
  Future<RangeGroup> getSparringPlacesRanges(int placeId);

  /// PUT /sparring_places/{id}/ranges
  ///
  /// Изменить серию. Участник имеет доступ только если это его место и серия активна. Если номера выстрелов в списке дублируются применяется последний встретившийся по порядку. При изменении организатором завершённой серии, результаты должны повлиять на открытие последующих серий и итоговый счёт (см. окончание серии).
  ///
  /// Important
  ///
  /// Изменение серии - транзакция, если во время изменения произошла ошибка, частичные изменения недопустимы
  ///
  /// Исключения:
  /// - NotFoundException(Место или серия не найдены)
  /// - InvalidParametersException(Неверные параметры)
  /// - InvalidScoreException(Счет выстрела не соответствует типу серии)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/sparring_places.md#put-sparring_placesidranges
  Future<Range> putSparringPlacesRange(int placeId, ChangeRange request);

  /// POST /sparring_places/{id}/ranges/{range_ordinal}/end
  ///
  /// Подтвердить, что серия закончена.
  ///
  /// 1. Завершить серию
  /// 2. Если счёт места ведётся по системе "общий счёт":
  ///    1. Если серия не последняя, активировать следующую серию
  ///    2. Иначе если противник уже закончил последнюю серию, проверить общий счёт:
  ///       1. Если общий счёт совпадает, создать перестрелки каждому участнику
  ///       2. Иначе перевести состояние спарринга в конечное
  /// 3. Если счёт места ведётся по системе "победные очки":
  ///    1. Если серия 1-ая или 2-ая, активировать следующую серию
  ///    2. Иначе если вы и противник уже закончили 3-ью или, последовательно, 4-ую серии, проверить общий счёт:
  ///       1. Если общий счёт совпадает, активировать следующую серию для обоих участников
  ///       2. Иначе перевести состояние спарринга в конечное
  ///    3. Иначе если вы и противник уже закончили 5-ую серию, проверить общий счёт:
  ///       1. Если общий счёт совпадает, создать перестрелки каждому участнику
  ///       2. Иначе перевести состояние спарринга в конечное
  ///
  /// Участник имеет доступ только если это его место.
  /// Серию можно завершить, если у всех выстрелов в серии выставлен non-null счёт.
  /// Завершение неактивной серии невозможно
  ///
  /// Исключения:
  /// - NotFoundException(Место или серия не найдены)
  /// - BadActionException(Невозможно завершить неполную или неактивную серию)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/sparring_places.md#post-sparring_placesidrangesrange_ordinalend
  Future<Range> endSparringPlacesRange(int placeId, int rangeOrdinal);

  /// PUT /sparring_places/{id}/shoot_out
  ///
  /// Изменить счёт перестрелки. Участник имеет доступ только если это его место и счёт перестрелки пуст.
  ///
  /// Исключения:
  /// - NotFoundException(Место или перестрелка не найдены)
  /// - InvalidParametersException(Неверные параметры)
  ///
  /// https://github.com/Glaz0k/archery-competitions-app/blob/feature/api-docs/docs/api/sparring_places.md#put-sparring_placesidshoot_out
  Future<ShootOut> putSparringPlacesShootOut(
    int placeId,
    ChangeShootOut request,
  );
}
