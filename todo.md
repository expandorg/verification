Domain model

job
  (id)
  verification_form
  task_form

  verification_module
  verification_agreement_count
  verification_score_threshold
  verification_prompt
  funding_verification_reward
  verification_limit

  task
    (job_id, id)
    data
    is_active
    assignment_count
    pending_count
    accepted_count
    verification_count

// task flow
assignment
  (job_id, task_id, id)
  user_id

responses
  (job_id, task_id, id)
  worker_id
  value
  is_accepted
  verifications_count

// verification flow
verification_assignment
  (job_id, task_id, response_id, id)
  user_id  

verification_responses
  (job_id, task_id, response_id, id)
  worker_id
  verifier_id
  value
  reason


