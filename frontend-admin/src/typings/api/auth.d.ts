declare namespace Api {
  /**
   * namespace Auth
   *
   * backend api module: "auth"
   */
  namespace Auth {
    interface LoginToken {
      token: string;
      refreshToken: string;
    }

    interface UserInfo {
      userId: string;
      userName: string;
      roles: string[];
      buttons: string[];
    }
  }

  /**
   * namespace Record
   *
   * backend api module: "records"
   */
  namespace Record {
    interface RecordItem {
      id: number;
      name: string;
      position: string;
      gender: string;
      title: string;
      content: string;
      images: string[];
      image_url: string;
      created_at: string;
      // Foreigner Profile
      date_of_birth: string;
      country_of_origin: string;
      passport_no: string;
      passport_issued_date: string;
      passport_expired_date: string;
      education_background: string;
      visa_entry_date: string;
      card_issued_date: string;
      card_expired_date: string;
      // Working History
      working_session: string;
      company_name: string;
      start_working_date: string;
      stop_working_date: string;
    }

    interface CreateRecordParams {
      name: string;
      position: string;
      gender: string;
      title: string;
      content: string;
      images: string[];
      // Foreigner Profile
      date_of_birth: string;
      country_of_origin: string;
      passport_no: string;
      passport_issued_date: string;
      passport_expired_date: string;
      education_background: string;
      visa_entry_date: string;
      card_issued_date: string;
      card_expired_date: string;
      // Working History
      working_session: string;
      company_name: string;
      start_working_date: string;
      stop_working_date: string;
    }
  }
}
